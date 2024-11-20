package processor

import (
	"encoding/json"
	"log"
	"personalNotificationService/factory"
	"personalNotificationService/kafka"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type NotificationProcessor struct {
	NotificationFactory factory.NotificationFactory
}

func NewNotificationProcessor(notificationFactory factory.NotificationFactory) kafka.MessageProcessor {
	return &NotificationProcessor{NotificationFactory: notificationFactory}
}

func (o *NotificationProcessor) ProcessMessage(message *sarama.ConsumerMessage) error {
	log.Printf("Processing notification message: %s\n", string(message.Value))

	msg := kafka.NotificationMessage{}
	if err := json.Unmarshal(message.Value, &msg); err != nil {
		log.Printf("Error unmarshaling of notification message: %v", string(message.Value))
		return err
	}

	if err := o.NotificationFactory.CreateNotification(msg.Username, uuid.MustParse(msg.NotificationID)); err != nil {
		return err
	}

	log.Printf("Successful processing of notification message: %v", msg)
	return nil
}
