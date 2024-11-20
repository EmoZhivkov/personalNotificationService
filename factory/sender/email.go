package sender

import (
	"encoding/json"
	"errors"
	"log"
	"personalNotificationService/common"
	"personalNotificationService/repositories"

	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	Config     common.Config
	MailClient common.MailClient
}

func newEmailSender(params Params) (senderConcreteImpl, error) {
	return &EmailSender{
		Config:     params.Config,
		MailClient: params.MailClient,
	}, nil
}

func (e *EmailSender) Send(userNotificationSettings any, content common.Content) error {
	jsonData, err := json.Marshal(userNotificationSettings)
	if err != nil {
		return err
	}

	var emailNotificationSettings repositories.EmailNotificationSettings
	if err := json.Unmarshal(jsonData, &emailNotificationSettings); err != nil {
		return errors.New("error invalid notification settings")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "senderConcreteImpl@example.com") // TODO: change this in the future
	m.SetHeader("To", emailNotificationSettings.UserEmail)
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/plain", string(content))

	_ = e.MailClient.DialAndSend(m)

	log.Print("Successfully sent email message")
	return nil
}
