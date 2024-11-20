package repositories

import "github.com/google/uuid"

type UserNotificationChannelsDatabase interface {
	CreateUserNotificationChannel(channel *UserNotificationChannels) error
	BulkCreateUserNotificationChannel(userNotificationChannels []UserNotificationChannels) error
	GetUserNotificationChannels(username string, notificationID uuid.UUID) (NotificationChannels, error)
}

type userNotificationChannelsDatabase struct {
	client DbClient
}

func NewUserNotificationChannelsDatabase(client DbClient) UserNotificationChannelsDatabase {
	return &userNotificationChannelsDatabase{client: client}
}

func (db *userNotificationChannelsDatabase) CreateUserNotificationChannel(channel *UserNotificationChannels) error {
	if err := db.client.Create(channel).Error; err != nil {
		return err
	}
	return nil
}

func (db *userNotificationChannelsDatabase) BulkCreateUserNotificationChannel(userNotificationChannels []UserNotificationChannels) error {
	if err := db.client.Create(&userNotificationChannels).Error; err != nil {
		return err
	}
	return nil
}

func (db *userNotificationChannelsDatabase) GetUserNotificationChannels(username string, notificationID uuid.UUID) (NotificationChannels, error) {
	var channels UserNotificationChannels
	if err := db.client.Take(&channels, "username = ? AND notification_id = ?", username, notificationID).Error; err != nil {
		return nil, err
	}
	return channels.Channels, nil
}
