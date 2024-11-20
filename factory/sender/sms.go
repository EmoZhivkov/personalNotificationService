package sender

import (
	"encoding/json"
	"errors"
	"log"
	"personalNotificationService/common"
	"personalNotificationService/repositories"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
)

type SmsSender struct {
	Config    common.Config
	SmsClient common.SmsClient
}

func newSmsSender(params Params) (senderConcreteImpl, error) {
	return &SmsSender{
		Config:    params.Config,
		SmsClient: params.SmsClient,
	}, nil
}

func (e *SmsSender) Send(userNotificationSettings any, content common.Content) error {
	jsonData, err := json.Marshal(userNotificationSettings)
	if err != nil {
		return err
	}

	var smsNotificationSettings repositories.SmsNotificationSettings
	if err := json.Unmarshal(jsonData, &smsNotificationSettings); err != nil {
		return errors.New("error invalid notification settings")
	}

	_, _ = e.SmsClient.Submit(&smpp.ShortMessage{
		Src:  "12345",
		Dst:  smsNotificationSettings.UserNumber,
		Text: pdutext.Raw(content),
	})

	log.Print("Successfully sent sms message")
	return err
}
