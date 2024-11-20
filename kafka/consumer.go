package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type MessageProcessor interface {
	ProcessMessage(message *sarama.ConsumerMessage) error
}

type Consumer struct {
	Brokers   []string
	Topic     string
	GroupID   string
	Processor MessageProcessor
}

func (kc *Consumer) StartConsuming() error {
	config := sarama.NewConfig()
	config.Consumer.Offsets.AutoCommit.Enable = false // in order to ensure "at least once" SLA
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(kc.Brokers, kc.GroupID, config)
	if err != nil {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}
	defer consumerGroup.Close()

	handler := ConsumerGroupHandler{
		Processor: kc.Processor,
		Topic:     kc.Topic,
	}

	for {
		err := consumerGroup.Consume(context.Background(), []string{kc.Topic}, &handler)
		if err != nil {
			log.Printf("Error consuming messages from topic %s: %v\n", kc.Topic, err)
		}
	}
}

type ConsumerGroupHandler struct {
	Processor MessageProcessor
	Topic     string
}

func (ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error { return nil }

func (ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	maxRetries := 3
	for message := range claim.Messages() {
		retryDelay := 2 * time.Second

		var err error
		for attempt := 1; attempt <= maxRetries; attempt++ {
			err = h.Processor.ProcessMessage(message)
			if err == nil {
				session.MarkMessage(message, "")
				session.Commit()
				break
			}

			log.Printf("Retrying message from topic %s (Attempt %d/%d): %v\n", h.Topic, attempt, maxRetries, err)

			if attempt < maxRetries {
				time.Sleep(retryDelay)
				retryDelay *= 2
			}
		}

		if err != nil {
			// TODO: in the future this could be handled with a Dead Letter Queue for further retry logic
			log.Printf("Error processing message from topic %s: %v\n", h.Topic, err)
			session.MarkMessage(message, "")
			session.Commit()
		}
	}
	return nil
}
