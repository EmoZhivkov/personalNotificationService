package sender

import (
	"encoding/json"
	"errors"
	"log"
	"personalNotificationService/common"
	"personalNotificationService/repositories"

	"github.com/slack-go/slack"
)

type SlackSender struct {
	Config      common.Config
	SlackClient common.SlackClient
}

func newSlackSender(params Params) (senderConcreteImpl, error) {
	return &SlackSender{
		Config:      params.Config,
		SlackClient: params.SlackClient,
	}, nil
}

func (e *SlackSender) Send(userNotificationSettings any, content common.Content) error {
	jsonData, err := json.Marshal(userNotificationSettings)
	if err != nil {
		return err
	}

	var slackNotificationSettings repositories.SlackNotificationSettings
	if err := json.Unmarshal(jsonData, &slackNotificationSettings); err != nil {
		return errors.New("error invalid notification settings")
	}

	_, _, _ = e.SlackClient.PostMessage(
		e.Config.SlackChannelID,
		slack.MsgOptionText(string(content), false),
	)

	log.Print("Successfully sent slack message")
	return nil
}
