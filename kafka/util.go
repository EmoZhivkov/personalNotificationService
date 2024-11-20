package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

func GetPartitions(config *sarama.Config, kafkaBrokerURL, topic string) ([]int32, error) {
	client, err := sarama.NewClient([]string{kafkaBrokerURL}, config)
	if err != nil {
		return nil, fmt.Errorf("error instantiating Kafka client: %w", err)
	}
	defer client.Close()

	partitions, err := client.Partitions(topic)
	if err != nil {
		return nil, fmt.Errorf("error getting partitions for topic %s: %w", topic, err)
	}

	return partitions, nil
}
