package factory

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"personalNotificationService/common"
	contentmocks "personalNotificationService/factory/content/mocks"
	metadatamocks "personalNotificationService/factory/metadata/mocks"
	sendermocks "personalNotificationService/factory/sender/mocks"
	"personalNotificationService/repositories"
	dbmocks "personalNotificationService/repositories/mocks"
)

func TestNotificationFactory_CreateNotification_Success(t *testing.T) {
	mockUserDB := dbmocks.NewMockUserDatabase(t)
	mockNotificationDB := dbmocks.NewMockNotificationDatabase(t)
	mockTemplateDB := dbmocks.NewMockTemplateDatabase(t)
	mockUserNotificationChannelsDB := dbmocks.NewMockUserNotificationChannelsDatabase(t)
	mockMetadataGenerator := metadatamocks.NewMockGenerator(t)
	mockContentGenerator := contentmocks.NewMockGenerator(t)
	mockNotificationSender := sendermocks.NewMockSender(t)

	targetUsername := "test_user"
	notificationID := uuid.New()
	templateID := uuid.New()

	mockNotificationDB.On("GetNotificationByID", notificationID).Return(&repositories.Notification{
		ID:   notificationID,
		Type: repositories.NotificationType("test"),
		ChannelToTemplateID: map[repositories.NotificationChannel]uuid.UUID{
			repositories.EmailNotificationChannel: templateID,
		},
	}, nil)

	mockMetadataGenerator.On("GenerateMetadata", repositories.NotificationType("test")).Return(common.Metadata{"key": "value"}, nil)

	mockUserNotificationChannelsDB.On("GetUserNotificationChannels", targetUsername, notificationID).Return(repositories.NotificationChannels{
		repositories.EmailNotificationChannel,
	}, nil)

	mockUserDB.On("GetUserByUsername", targetUsername).Return(&repositories.User{
		NotificationSettings: repositories.NotificationSettings{
			repositories.EmailNotificationChannel: repositories.EmailNotificationSettings{UserEmail: "test@example.com"},
		},
	}, nil)

	mockTemplateDB.On("GetTemplateByID", templateID).Return(&repositories.Template{
		ID:       templateID,
		Template: "Hello, {{.key}}",
	}, nil)

	mockContentGenerator.On("GenerateContent", common.Metadata{"key": "value"}, mock.Anything).Return(common.Content("Hello, value"), nil)

	mockNotificationSender.On("Send", repositories.EmailNotificationChannel, repositories.EmailNotificationSettings{UserEmail: "test@example.com"}, common.Content("Hello, value")).Return(nil)

	factory := NewNotificationFactory(Params{
		UserDB:                     mockUserDB,
		NotificationDB:             mockNotificationDB,
		TemplateDB:                 mockTemplateDB,
		UserNotificationChannelsDB: mockUserNotificationChannelsDB,
		MetadataGenerator:          mockMetadataGenerator,
		ContentGenerator:           mockContentGenerator,
		NotificationSender:         mockNotificationSender,
	})

	err := factory.CreateNotification(targetUsername, notificationID)
	assert.NoError(t, err)
}

func TestNotificationFactory_CreateNotification_NotificationNotFound(t *testing.T) {
	mockNotificationDB := dbmocks.NewMockNotificationDatabase(t)
	targetUsername := "test_user"
	notificationID := uuid.New()

	mockNotificationDB.On("GetNotificationByID", notificationID).Return(nil, errors.New("notification not found"))

	factory := NewNotificationFactory(Params{
		NotificationDB: mockNotificationDB,
	})

	err := factory.CreateNotification(targetUsername, notificationID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "notification not found")
}

func TestNotificationFactory_CreateNotification_ChannelsNotConfigured(t *testing.T) {
	mockUserDB := dbmocks.NewMockUserDatabase(t)
	mockNotificationDB := dbmocks.NewMockNotificationDatabase(t)
	mockTemplateDB := dbmocks.NewMockTemplateDatabase(t)
	mockUserNotificationChannelsDB := dbmocks.NewMockUserNotificationChannelsDatabase(t)
	mockMetadataGenerator := metadatamocks.NewMockGenerator(t)
	mockContentGenerator := contentmocks.NewMockGenerator(t)
	mockNotificationSender := sendermocks.NewMockSender(t)

	targetUsername := "test_user"
	notificationID := uuid.New()

	mockNotificationDB.On("GetNotificationByID", notificationID).Return(&repositories.Notification{
		ID:   notificationID,
		Type: repositories.NotificationType("test"),
	}, nil)

	mockMetadataGenerator.On("GenerateMetadata", repositories.NotificationType("test")).Return(common.Metadata{"key": "value"}, nil)

	mockUserNotificationChannelsDB.On("GetUserNotificationChannels", targetUsername, notificationID).Return(
		nil,
		errors.New("channels not configured"),
	)

	factory := NewNotificationFactory(Params{
		UserDB:                     mockUserDB,
		NotificationDB:             mockNotificationDB,
		TemplateDB:                 mockTemplateDB,
		UserNotificationChannelsDB: mockUserNotificationChannelsDB,
		MetadataGenerator:          mockMetadataGenerator,
		ContentGenerator:           mockContentGenerator,
		NotificationSender:         mockNotificationSender,
	})
	err := factory.CreateNotification(targetUsername, notificationID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "channels not configured")
}

func TestNotificationFactory_CreateNotification_SenderFailure(t *testing.T) {
	mockUserDB := dbmocks.NewMockUserDatabase(t)
	mockNotificationDB := dbmocks.NewMockNotificationDatabase(t)
	mockTemplateDB := dbmocks.NewMockTemplateDatabase(t)
	mockUserNotificationChannelsDB := dbmocks.NewMockUserNotificationChannelsDatabase(t)
	mockMetadataGenerator := metadatamocks.NewMockGenerator(t)
	mockContentGenerator := contentmocks.NewMockGenerator(t)
	mockNotificationSender := sendermocks.NewMockSender(t)

	targetUsername := "test_user"
	notificationID := uuid.New()
	templateID := uuid.New()

	mockNotificationDB.On("GetNotificationByID", notificationID).Return(&repositories.Notification{
		ID:   notificationID,
		Type: repositories.NotificationType("test"),
		ChannelToTemplateID: map[repositories.NotificationChannel]uuid.UUID{
			repositories.EmailNotificationChannel: templateID,
		},
	}, nil)

	mockMetadataGenerator.On("GenerateMetadata", repositories.NotificationType("test")).Return(common.Metadata{"key": "value"}, nil)

	mockUserNotificationChannelsDB.On("GetUserNotificationChannels", targetUsername, notificationID).Return(repositories.NotificationChannels{
		repositories.EmailNotificationChannel,
	}, nil)

	mockUserDB.On("GetUserByUsername", targetUsername).Return(&repositories.User{
		NotificationSettings: repositories.NotificationSettings{
			repositories.EmailNotificationChannel: repositories.EmailNotificationSettings{UserEmail: "test@example.com"},
		},
	}, nil)

	mockTemplateDB.On("GetTemplateByID", templateID).Return(&repositories.Template{
		ID:       templateID,
		Template: "Hello, {{.key}}",
	}, nil)

	mockContentGenerator.On("GenerateContent", common.Metadata{"key": "value"}, mock.Anything).Return(common.Content("Hello, value"), nil)

	mockNotificationSender.On("Send", repositories.EmailNotificationChannel, repositories.EmailNotificationSettings{UserEmail: "test@example.com"}, common.Content("Hello, value")).Return(errors.New("failed to send email"))

	factory := NewNotificationFactory(Params{
		UserDB:                     mockUserDB,
		NotificationDB:             mockNotificationDB,
		TemplateDB:                 mockTemplateDB,
		UserNotificationChannelsDB: mockUserNotificationChannelsDB,
		MetadataGenerator:          mockMetadataGenerator,
		ContentGenerator:           mockContentGenerator,
		NotificationSender:         mockNotificationSender,
	})

	err := factory.CreateNotification(targetUsername, notificationID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send email")
}
