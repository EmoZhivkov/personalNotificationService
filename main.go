package main

import (
	"log"
	"personalNotificationService/kafka"
)

func main() {
	producer, err := kafka.NewPriorityProducer("localhost:9093", kafka.NotificationsTopic)
	if err != nil {
		log.Fatal("error instantiating producer")
	}

	for i := 0; i <= 10; i++ {
		if err := producer.SendMessage(&kafka.NotificationMessage{
			Username:       "alice",
			NotificationID: "4f343b21-70ad-48bc-9557-bdbfb477136d",
			Priority:       kafka.HighMessagePriority,
		}); err != nil {
			log.Fatalf("error sending router msg %v", err)
		}
	}

	//for i := 0; i <= 10; i++ {
	//	if err := producer.SendMessage(&kafka.NotificationMessage{
	//		Username:         "gosho",
	//		NotificationType: "",
	//		Priority:         kafka.MediumMessagePriority,
	//	}); err != nil {
	//		log.Fatalf("error sending router msg %v", err)
	//	}
	//}
	//
	//for i := 0; i <= 10; i++ {
	//	if err := producer.SendMessage(&kafka.NotificationMessage{
	//		Username:         "pencho",
	//		NotificationType: "",
	//		Priority:         kafka.LowMessagePriority,
	//	}); err != nil {
	//		log.Fatalf("error sending router msg %v", err)
	//	}
	//}
}
