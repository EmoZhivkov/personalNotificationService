package kafka

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/IBM/sarama"
)

type MessageWithPriority interface {
	GetPriority() MessagePriority
	Encode() ([]byte, error)
}

type PriorityProducer interface {
	SendMessage(MessageWithPriority) error
	Close() error
}

type priorityProducer struct {
	producer           sarama.SyncProducer
	Topic              string
	PriorityPartitions map[MessagePriority][]int32
}

func NewPriorityProducer(kafkaBrokerURL string, topic string) (PriorityProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // in order to wait for an acknowledgment that the message is sent
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewManualPartitioner // in order to use partitioning logic based on priority

	syncProducer, err := sarama.NewSyncProducer([]string{kafkaBrokerURL}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka priorityProducer: %w", err)
	}

	partitions, err := GetPartitions(config, kafkaBrokerURL, topic)
	if err != nil {
		return nil, err
	}

	priorityPartitions, err := GeneratePriorityPartitions(partitions, PriorityToRatio)
	if err != nil {
		return nil, err
	}

	return &priorityProducer{
		producer:           syncProducer,
		Topic:              topic,
		PriorityPartitions: priorityPartitions,
	}, nil
}

func (kp *priorityProducer) SendMessage(message MessageWithPriority) error {
	partitions, ok := kp.PriorityPartitions[message.GetPriority()]
	if !ok || len(partitions) == 0 {
		return fmt.Errorf("invalid or unassigned priority level: %s", message.GetPriority())
	}

	partition := partitions[rand.Intn(len(partitions))]

	encodedValue, err := message.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic:     kp.Topic,
		Value:     sarama.ByteEncoder(encodedValue),
		Partition: partition,
	}

	actualPartition, offset, err := kp.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	if actualPartition != partition {
		log.Printf("Warning: Message was routed to partition %d instead of %d", actualPartition, partition)
	}

	log.Printf("Message sent to topic %s: partition=%d, offset=%d\n", kp.Topic, partition, offset)
	return nil
}

func (kp *priorityProducer) Close() error {
	return kp.producer.Close()
}
