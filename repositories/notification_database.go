package repositories

import "github.com/google/uuid"

type NotificationDatabase interface {
	CreateNotification(notification *Notification) error
	CreateNotifications(notifications Notifications) error
	GetNotificationByID(id uuid.UUID) (*Notification, error)
	GetNotificationsByIDs(ids uuid.UUIDs) (Notifications, error)
	GetNotificationsByType(notificationType NotificationType) ([]Notification, error)
}

type notificationDatabase struct {
	client DbClient
}

func NewNotificationDatabase(client DbClient) NotificationDatabase {
	return &notificationDatabase{client: client}
}

func (db *notificationDatabase) CreateNotification(notification *Notification) error {
	if err := db.client.Create(notification).Error; err != nil {
		return err
	}
	return nil
}

func (db *notificationDatabase) CreateNotifications(notifications Notifications) error {
	if err := db.client.Create(&notifications).Error; err != nil {
		return err
	}
	return nil
}

func (db *notificationDatabase) GetNotificationByID(id uuid.UUID) (*Notification, error) {
	var notification Notification
	if err := db.client.Take(&notification, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}

func (db *notificationDatabase) GetNotificationsByIDs(ids uuid.UUIDs) (Notifications, error) {
	var notifications Notifications
	if err := db.client.Find(&notifications, "id IN ?", ids).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

func (db *notificationDatabase) GetNotificationsByType(notificationType NotificationType) ([]Notification, error) {
	var notifications []Notification
	if err := db.client.Find(&notifications, "type = ?", notificationType).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}
