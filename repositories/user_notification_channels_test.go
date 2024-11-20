package repositories

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserNotificationChannel(t *testing.T) {
	testUserNotifDB := NewUserNotificationChannelsDatabase(getTestInMemoryDBClient())

	channel := &UserNotificationChannels{
		Username:       "test_user",
		NotificationID: uuid.New(),
		Channels:       NotificationChannels{EmailNotificationChannel, SmsNotificationChannel},
	}
	assert.NoError(t, testUserNotifDB.CreateUserNotificationChannel(channel))

	retrievedChannels, err := testUserNotifDB.GetUserNotificationChannels(channel.Username, channel.NotificationID)
	assert.NoError(t, err)
	assert.Equal(t, channel.Channels, retrievedChannels)
}

func TestBulkCreateUserNotificationChannels(t *testing.T) {
	testUserNotifDB := NewUserNotificationChannelsDatabase(getTestInMemoryDBClient())

	channels := []UserNotificationChannels{
		{
			Username:       "user1",
			NotificationID: uuid.New(),
			Channels:       NotificationChannels{EmailNotificationChannel},
		},
		{
			Username:       "user2",
			NotificationID: uuid.New(),
			Channels:       NotificationChannels{SlackNotificationChannel},
		},
	}

	assert.NoError(t, testUserNotifDB.BulkCreateUserNotificationChannel(channels))

	for _, channel := range channels {
		retrievedChannels, err := testUserNotifDB.GetUserNotificationChannels(channel.Username, channel.NotificationID)
		assert.NoError(t, err)
		assert.Equal(t, channel.Channels, retrievedChannels)
	}
}

func TestGetUserNotificationChannels(t *testing.T) {
	testUserNotifDB := NewUserNotificationChannelsDatabase(getTestInMemoryDBClient())

	_, err := testUserNotifDB.GetUserNotificationChannels("non_existent_user", uuid.New())
	assert.Error(t, err)

	channel := &UserNotificationChannels{
		Username:       "test_user",
		NotificationID: uuid.New(),
		Channels:       NotificationChannels{SmsNotificationChannel},
	}

	assert.NoError(t, testUserNotifDB.CreateUserNotificationChannel(channel))

	retrievedChannels, err := testUserNotifDB.GetUserNotificationChannels(channel.Username, channel.NotificationID)
	assert.NoError(t, err)
	assert.Equal(t, channel.Channels, retrievedChannels)
}
