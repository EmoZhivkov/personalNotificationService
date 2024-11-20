package factory

import (
	"fmt"
	"log"
	"personalNotificationService/common"
	"personalNotificationService/factory/content"
	"personalNotificationService/factory/metadata"
	"personalNotificationService/factory/sender"
	"personalNotificationService/repositories"

	"github.com/google/uuid"
)

type NotificationFactory interface {
	CreateNotification(targetUsername string, notificationID uuid.UUID) error
}

type Params struct {
	Config                     common.Config
	UserDB                     repositories.UserDatabase
	NotificationDB             repositories.NotificationDatabase
	TemplateDB                 repositories.TemplateDatabase
	UserNotificationChannelsDB repositories.UserNotificationChannelsDatabase
	MetadataGenerator          metadata.Generator
	ContentGenerator           content.Generator
	NotificationSender         sender.Sender
}

type notificationFactory struct {
	Params
}

func NewNotificationFactory(params Params) NotificationFactory {
	return &notificationFactory{Params: params}
}

func (n *notificationFactory) CreateNotification(targetUsername string, notificationID uuid.UUID) error {
	notification, err := n.NotificationDB.GetNotificationByID(notificationID)
	if err != nil {
		return err
	}

	notificationMetadata, err := n.MetadataGenerator.GenerateMetadata(notification.Type)
	if err != nil {
		return err
	}

	userNotificationChannelSet, err := n.getConfiguredUserNotificationChannels(targetUsername, notification)
	if err != nil {
		return err
	}

	userNotificationSettings, err := n.getUserNotificationSettings(targetUsername)
	if err != nil {
		return err
	}

	for notificationChannel, templateID := range notification.ChannelToTemplateID {
		if _, exists := userNotificationChannelSet[notificationChannel]; !exists {
			continue // user has not configured this channel for this notification type
		}

		notificationTemplate, err := n.TemplateDB.GetTemplateByID(templateID)
		if err != nil {
			return err
		}

		notificationContent, err := n.ContentGenerator.GenerateContent(notificationMetadata, *notificationTemplate)
		if err != nil {
			return err
		}

		notificationSettings, ok := userNotificationSettings[notificationChannel]
		if !ok {
			return fmt.Errorf("error not existing notification settings for channel: %v", notificationChannel)
		}

		if err := n.NotificationSender.Send(notificationChannel, notificationSettings, notificationContent); err != nil {
			log.Printf("error failed to send notification on channel: %v", notificationChannel)
			return err
		}
	}
	return nil
}

func (n *notificationFactory) getConfiguredUserNotificationChannels(targetUsername string, notification *repositories.Notification) (map[repositories.NotificationChannel]struct{}, error) {
	userNotificationChannels, err := n.UserNotificationChannelsDB.GetUserNotificationChannels(targetUsername, notification.ID)
	if err != nil {
		return nil, err
	}
	return userNotificationChannels.ToNotificationChannelSet(), nil
}

func (n *notificationFactory) getUserNotificationSettings(targetUsername string) (repositories.NotificationSettings, error) {
	user, err := n.UserDB.GetUserByUsername(targetUsername)
	if err != nil {
		return nil, err
	}
	return user.NotificationSettings, err
}
