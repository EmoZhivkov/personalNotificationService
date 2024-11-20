package sender

import (
	"errors"
	"personalNotificationService/common"
	"personalNotificationService/repositories"
)

type senderConcreteImpl interface {
	Send(any, common.Content) error
}

var senderFactoryMethods = map[repositories.NotificationChannel]func(Params) (senderConcreteImpl, error){
	repositories.EmailNotificationChannel: newEmailSender,
	repositories.SlackNotificationChannel: newSlackSender,
	repositories.SmsNotificationChannel:   newSmsSender,
}

type Sender interface {
	Send(repositories.NotificationChannel, any, common.Content) error
}

func NewNotificationSender(params Params) Sender {
	return &sender{Params: params}
}

type Params struct {
	Config      common.Config
	MailClient  common.MailClient
	SlackClient common.SlackClient
	SmsClient   common.SmsClient
}

type sender struct {
	Params
}

func (s *sender) Send(notificationChannel repositories.NotificationChannel, notificationSettings any, content common.Content) error {
	factoryMethod, ok := senderFactoryMethods[notificationChannel]
	if !ok {
		return errors.New("error invalid notification channel")
	}

	senderImpl, err := factoryMethod(s.Params)
	if err != nil {
		return errors.New("error instantiating sender concrete implementation")
	}

	return senderImpl.Send(notificationSettings, content)
}
