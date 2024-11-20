package repositories

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateNotification(t *testing.T) {
	testNotificationDB := NewNotificationDatabase(getTestInMemoryDBClient())

	notification := &Notification{
		ID:       uuid.New(),
		Type:     SuccessfulTransactionNotificationType,
		Priority: HighNotificationPriority,
		ChannelToTemplateID: ChannelToTemplateMap{
			EmailNotificationChannel: uuid.New(),
			SmsNotificationChannel:   uuid.New(),
		},
	}

	assert.NoError(t, testNotificationDB.CreateNotification(notification))

	retrievedNotification, err := testNotificationDB.GetNotificationByID(notification.ID)
	assert.NoError(t, err)
	assert.Equal(t, notification.ID, retrievedNotification.ID)
	assert.Equal(t, notification.Type, retrievedNotification.Type)
	assert.Equal(t, notification.Priority, retrievedNotification.Priority)
	assert.Equal(t, notification.ChannelToTemplateID, retrievedNotification.ChannelToTemplateID)
}

func TestCreateNotifications(t *testing.T) {
	testNotificationDB := NewNotificationDatabase(getTestInMemoryDBClient())

	notifications := Notifications{
		{
			ID:       uuid.New(),
			Type:     FailedTransactionNotificationType,
			Priority: MediumNotificationPriority,
			ChannelToTemplateID: ChannelToTemplateMap{
				EmailNotificationChannel: uuid.New(),
			},
		},
		{
			ID:       uuid.New(),
			Type:     AmountReceivedNotificationType,
			Priority: LowNotificationPriority,
			ChannelToTemplateID: ChannelToTemplateMap{
				SlackNotificationChannel: uuid.New(),
			},
		},
	}
	assert.NoError(t, testNotificationDB.CreateNotifications(notifications))

	notification1, err := testNotificationDB.GetNotificationByID(notifications[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, notifications[0].ID, notification1.ID)
	assert.Equal(t, notifications[0].Type, notification1.Type)
	assert.Equal(t, notifications[0].Priority, notification1.Priority)
	assert.Equal(t, notifications[0].ChannelToTemplateID, notification1.ChannelToTemplateID)

	notification2, err := testNotificationDB.GetNotificationByID(notifications[1].ID)
	assert.NoError(t, err)
	assert.Equal(t, notifications[1].ID, notification2.ID)
	assert.Equal(t, notifications[1].Type, notification2.Type)
	assert.Equal(t, notifications[1].Priority, notification2.Priority)
	assert.Equal(t, notifications[1].ChannelToTemplateID, notification2.ChannelToTemplateID)
}

func TestGetNotificationByID(t *testing.T) {
	testNotificationDB := NewNotificationDatabase(getTestInMemoryDBClient())

	nonExistentID := uuid.New()
	_, err := testNotificationDB.GetNotificationByID(nonExistentID)
	assert.Error(t, err)
}

func TestGetNotificationsByType(t *testing.T) {
	testNotificationDB := NewNotificationDatabase(getTestInMemoryDBClient())

	notifications := []Notification{
		{
			ID:       uuid.New(),
			Type:     SuccessfulTransactionNotificationType,
			Priority: HighNotificationPriority,
		},
		{
			ID:       uuid.New(),
			Type:     FailedTransactionNotificationType,
			Priority: MediumNotificationPriority,
		},
		{
			ID:       uuid.New(),
			Type:     SuccessfulTransactionNotificationType,
			Priority: LowNotificationPriority,
		},
	}

	for _, n := range notifications {
		assert.NoError(t, testNotificationDB.CreateNotification(&n))
	}

	successfulNotifications, err := testNotificationDB.GetNotificationsByType(SuccessfulTransactionNotificationType)
	assert.NoError(t, err)
	assert.Len(t, successfulNotifications, 2)

	for _, n := range successfulNotifications {
		assert.Equal(t, SuccessfulTransactionNotificationType, n.Type)
	}

	failedNotifications, err := testNotificationDB.GetNotificationsByType(FailedTransactionNotificationType)
	assert.NoError(t, err)
	assert.Len(t, failedNotifications, 1)

	assert.Equal(t, FailedTransactionNotificationType, failedNotifications[0].Type)
}
