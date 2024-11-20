package main

import (
	"log"
	"personalNotificationService/common"
	"personalNotificationService/factory"
	"personalNotificationService/factory/content"
	"personalNotificationService/factory/metadata"
	"personalNotificationService/factory/sender"
	"personalNotificationService/kafka"
	"personalNotificationService/processor"
	"personalNotificationService/repositories"
	"strconv"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"gopkg.in/gomail.v2"
)

func main() {
	log.Println("Starting Notification Consumer")

	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}
	config := common.NewConfig()
	log.Println("Loaded configuration")

	client, err := repositories.NewDbClient(config.DbConnectionURL)
	if err != nil {
		log.Fatalf("Error creating Database client: %v", err)
	}

	userDB := repositories.NewUserDatabase(client)
	notificationDB := repositories.NewNotificationDatabase(client)
	templateDB := repositories.NewTemplateDatabase(client)
	userNotificationChannelsDB := repositories.NewUserNotificationChannelsDatabase(client)

	smtpPort, err := strconv.Atoi(config.SmtpPort)
	if err != nil {
		log.Fatal("invalid SMTP_PORT")
	}
	mailClient := gomail.NewDialer(config.SmtpHost, smtpPort, config.SmtpUsername, config.SmtpPassword)

	slackClient := slack.New(config.SlackBotToken)

	smsClient := &smpp.Transceiver{
		Addr:   config.SmscHost + ":" + config.SmscPort,
		User:   config.SmscUsername,
		Passwd: config.SmscPassword,
	}

	connCh := smsClient.Bind()
	defer smsClient.Close()

	go func() {
		select {
		case connStatus := <-connCh:
			if connStatus.Error() != nil {
				log.Fatalf("Failed to bind to SMSC: %v", err)
			}
		}
	}()

	metadataGenerator := metadata.NewMetadataGenerator(metadata.Params{})
	contentGenerator := content.NewContentGenerator(content.Params{})
	notificationSender := sender.NewNotificationSender(sender.Params{
		Config:      config,
		MailClient:  mailClient,
		SlackClient: slackClient,
		SmsClient:   smsClient,
	})

	notificationFacory := factory.NewNotificationFactory(factory.Params{
		Config:                     config,
		UserDB:                     userDB,
		NotificationDB:             notificationDB,
		TemplateDB:                 templateDB,
		UserNotificationChannelsDB: userNotificationChannelsDB,
		MetadataGenerator:          metadataGenerator,
		ContentGenerator:           contentGenerator,
		NotificationSender:         notificationSender,
	})

	emailConsumer := kafka.Consumer{
		Brokers:   []string{config.KafkaBrokerURL},
		Topic:     kafka.NotificationsTopic,
		GroupID:   kafka.NotificationsTopic + "-group",
		Processor: processor.NewNotificationProcessor(notificationFacory),
	}

	if err := emailConsumer.StartConsuming(); err != nil {
		log.Fatalf("Notification Consumer failed: %v", err)
	}
}
