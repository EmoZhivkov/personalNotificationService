package main

import (
	"log"
	"personalNotificationService/api"
	"personalNotificationService/common"
	"personalNotificationService/kafka"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting API")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	config := common.NewConfig()

	producer, err := kafka.NewPriorityProducer(config.KafkaBrokerURL, kafka.NotificationsTopic)
	if err != nil {
		log.Fatal("error instantiating producer")
	}
	defer producer.Close()

	server := api.NewServer(config, producer)

	if err := server.Run(); err != nil {
		log.Fatalf("error running the server: %v", err.Error())
	}
}
