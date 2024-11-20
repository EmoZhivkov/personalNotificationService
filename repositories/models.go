package repositories

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type UserNotificationChannels struct {
	Username       string               `gorm:"column:username;foreignKey"`
	NotificationID uuid.UUID            `gorm:"column:notification_id;type:uuid;foreignKey"`
	Channels       NotificationChannels `gorm:"column:channels;type:jsonb" json:"channels"`
}

type Templates []Template

type Template struct {
	ID       uuid.UUID           `gorm:"column:id;type:uuid;primaryKey"`
	Channel  NotificationChannel `gorm:"column:channel"`
	Template string              `gorm:"column:template"`
}

type NotificationChannels []NotificationChannel

func (n NotificationChannels) ToNotificationChannelSet() map[NotificationChannel]struct{} {
	notificationChannelSet := make(map[NotificationChannel]struct{}, len(n))
	for _, notificationChannel := range n {
		tmp := notificationChannel
		notificationChannelSet[tmp] = struct{}{}
	}
	return notificationChannelSet
}

func (n NotificationChannels) Value() (driver.Value, error) {
	return json.Marshal(n)
}

func (n *NotificationChannels) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to convert database value to NotificationChannelArray: %v", value)
	}
	return json.Unmarshal(bytes, n)
}

type NotificationChannel string

func (n *NotificationChannel) Scan(value interface{}) error {
	*n = NotificationChannel(value.(string))
	return nil
}

func (n NotificationChannel) Value() (driver.Value, error) {
	return string(n), nil
}

const (
	EmailNotificationChannel NotificationChannel = "email"
	SlackNotificationChannel NotificationChannel = "slack"
	SmsNotificationChannel   NotificationChannel = "sms"
)

func IsValidNotificationChannel(channel string) bool {
	switch NotificationChannel(channel) {
	case EmailNotificationChannel, SlackNotificationChannel, SmsNotificationChannel:
		return true
	}
	return false
}

type ChannelToTemplateMap map[NotificationChannel]uuid.UUID

func (m ChannelToTemplateMap) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *ChannelToTemplateMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to convert database value to ChannelToTemplateMap: %v", value)
	}
	return json.Unmarshal(bytes, m)
}

type Notifications []Notification

type Notification struct {
	ID                  uuid.UUID            `gorm:"column:id;type:uuid;primaryKey"`
	Type                NotificationType     `gorm:"column:type"`
	Priority            NotificationPriority `gorm:"column:priority"`
	ChannelToTemplateID ChannelToTemplateMap `gorm:"column:channel_to_template;type:jsonb" json:"channel_to_template"`
}

type NotificationType string

func (n *NotificationType) Scan(value interface{}) error {
	*n = NotificationType(value.(string))
	return nil
}

func (n NotificationType) Value() (driver.Value, error) {
	return string(n), nil
}

const (
	SuccessfulTransactionNotificationType NotificationType = "successful_transaction"
	FailedTransactionNotificationType     NotificationType = "failed_transaction"
	AmountReceivedNotificationType        NotificationType = "amount_received"
)

func IsValidNotificationType(notificationType string) bool {
	switch NotificationType(notificationType) {
	case SuccessfulTransactionNotificationType, FailedTransactionNotificationType, AmountReceivedNotificationType:
		return true
	}
	return false
}

type NotificationPriority string

func (n *NotificationPriority) Scan(value interface{}) error {
	*n = NotificationPriority(value.(string))
	return nil
}

func (n NotificationPriority) Value() (driver.Value, error) {
	return string(n), nil
}

const (
	HighNotificationPriority   NotificationPriority = "high"
	MediumNotificationPriority NotificationPriority = "medium"
	LowNotificationPriority    NotificationPriority = "low"
)

func IsValidNotificationPriority(priority string) bool {
	switch NotificationPriority(priority) {
	case HighNotificationPriority, MediumNotificationPriority, LowNotificationPriority:
		return true
	}
	return false
}

type Users []*User

func (u Users) ToUserSet() map[string]*User {
	userSet := make(map[string]*User, len(u))
	for _, user := range u {
		tmp := user
		userSet[user.Username] = tmp
	}
	return userSet
}

type NotificationSettings map[NotificationChannel]interface{}

func (n NotificationSettings) Value() (driver.Value, error) {
	return json.Marshal(n)
}

func (n *NotificationSettings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to convert database value to NotificationSettings: %v", value)
	}
	return json.Unmarshal(bytes, n)
}

type User struct {
	Username             string               `gorm:"column:username;primaryKey"`
	PasswordHash         string               `gorm:"column:password_hash"`
	NotificationSettings NotificationSettings `gorm:"column:notification_settings;type:jsonb" json:"notification_settings"`
}

type EmailNotificationSettings struct {
	UserEmail string `json:"userEmail"`
}

type SlackNotificationSettings struct {
	UserHandle string `json:"userHandle"`
}

type SmsNotificationSettings struct {
	UserNumber string `json:"userNumber"`
}
