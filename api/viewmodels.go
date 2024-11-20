package api

import (
	"errors"
	"personalNotificationService/kafka"
	"personalNotificationService/repositories"

	"github.com/google/uuid"
)

type SendNotificationRequest struct {
	Username       string `json:"username"`
	NotificationID string `json:"notificationID"`
}

func (req *SendNotificationRequest) ToKafkaNotificationMessage(notification *repositories.Notification) *kafka.NotificationMessage {
	return &kafka.NotificationMessage{
		Username:       req.Username,
		NotificationID: req.NotificationID,
		Priority:       kafka.MessagePriority(notification.Priority),
	}
}

type CreateTemplateRequest struct {
	Channel  string `json:"channel"`
	Template string `json:"template"`
}

func (req *CreateTemplateRequest) Validate() error {
	if !repositories.IsValidNotificationChannel(req.Channel) {
		return errors.New("invalid notification channel")
	}
	return nil
}

func (req *CreateTemplateRequest) ToTemplateModel() repositories.Template {
	return repositories.Template{
		ID:       uuid.New(),
		Channel:  repositories.NotificationChannel(req.Channel),
		Template: req.Template,
	}
}

type TemplateResponse struct {
	ID       string `json:"id"`
	Channel  string `json:"channel"`
	Template string `json:"template"`
}

func ToTemplateResponse(template *repositories.Template) TemplateResponse {
	return TemplateResponse{
		ID:       template.ID.String(),
		Channel:  string(template.Channel),
		Template: template.Template,
	}
}

type NotificationRequest struct {
	Type                string              `json:"type"`
	Priority            string              `json:"priority"`
	ChannelToTemplateID []ChannelToTemplate `json:"channelToTemplateID"`
}

func (req *NotificationRequest) Validate() error {
	if !repositories.IsValidNotificationType(req.Type) {
		return errors.New("invalid notification type")
	}
	if !repositories.IsValidNotificationPriority(req.Priority) {
		return errors.New("invalid notification priority")
	}
	for _, channelToTemplate := range req.ChannelToTemplateID {
		if !repositories.IsValidNotificationChannel(channelToTemplate.Channel) {
			return errors.New("invalid notification channel")
		}
		if _, err := uuid.Parse(channelToTemplate.TemplateID); err != nil {
			return errors.New("invalid template ID")
		}
	}
	return nil
}

func (req *NotificationRequest) ToNotificationModel() repositories.Notification {
	channelToTemplateMap := make(map[repositories.NotificationChannel]uuid.UUID, len(req.ChannelToTemplateID))
	for _, channelToTemplate := range req.ChannelToTemplateID {
		channelToTemplateMap[repositories.NotificationChannel(channelToTemplate.Channel)] = uuid.MustParse(channelToTemplate.TemplateID)
	}

	return repositories.Notification{
		ID:                  uuid.New(),
		Type:                repositories.NotificationType(req.Type),
		Priority:            repositories.NotificationPriority(req.Priority),
		ChannelToTemplateID: channelToTemplateMap,
	}
}

type NotificationResponse struct {
	ID                  string              `json:"id"`
	Type                string              `json:"type"`
	Priority            string              `json:"priority"`
	ChannelToTemplateID []ChannelToTemplate `json:"channelToTemplateID"`
}

func ToNotificationResponse(notification *repositories.Notification) NotificationResponse {
	channelToTemplate := make([]ChannelToTemplate, 0, len(notification.ChannelToTemplateID))
	for channel, templateID := range notification.ChannelToTemplateID {
		channelToTemplate = append(channelToTemplate, ChannelToTemplate{
			Channel:    string(channel),
			TemplateID: templateID.String(),
		})
	}

	return NotificationResponse{
		ID:                  notification.ID.String(),
		Type:                string(notification.Type),
		Priority:            string(notification.Priority),
		ChannelToTemplateID: channelToTemplate,
	}
}

type ChannelToTemplate struct {
	Channel    string `json:"channel"`
	TemplateID string `json:"templateID"`
}

type UserNotificationChannelsRequest struct {
	Username       string   `json:"username"`
	NotificationID string   `json:"notificationID"`
	Channels       []string `json:"channels"`
}

func (req *UserNotificationChannelsRequest) Validate() error {
	if _, err := uuid.Parse(req.NotificationID); err != nil {
		return errors.New("invalid notification ID")
	}
	for _, channel := range req.Channels {
		if !repositories.IsValidNotificationChannel(channel) {
			return errors.New("invalid notification channel")
		}
	}
	return nil
}

func (req *UserNotificationChannelsRequest) ToUserNotificationChannelsModel() repositories.UserNotificationChannels {
	channels := make(repositories.NotificationChannels, 0, len(req.Channels))
	for _, channel := range req.Channels {
		channels = append(channels, repositories.NotificationChannel(channel))
	}

	return repositories.UserNotificationChannels{
		Username:       req.Username,
		NotificationID: uuid.MustParse(req.NotificationID),
		Channels:       channels,
	}
}

type UserNotificationChannelsResponse struct {
	Username       string   `json:"username"`
	NotificationID string   `json:"notificationID"`
	Channels       []string `json:"channels"`
}

func ToUserNotificationChannelsResponse(userNotificationChannels *repositories.UserNotificationChannels) UserNotificationChannelsResponse {
	channels := make([]string, 0, len(userNotificationChannels.Channels))
	for _, channel := range userNotificationChannels.Channels {
		channels = append(channels, string(channel))
	}

	return UserNotificationChannelsResponse{
		Username:       userNotificationChannels.Username,
		NotificationID: userNotificationChannels.NotificationID.String(),
		Channels:       channels,
	}
}
