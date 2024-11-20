package common

import (
	"github.com/fiorix/go-smpp/smpp"
	"github.com/slack-go/slack"
	"gopkg.in/gomail.v2"
)

type MailClient interface {
	DialAndSend(m ...*gomail.Message) error
}

type SlackClient interface {
	PostMessage(channelID string, options ...slack.MsgOption) (string, string, error)
}

type SmsClient interface {
	Submit(sm *smpp.ShortMessage) (*smpp.ShortMessage, error)
}
