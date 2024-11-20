package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"personalNotificationService/kafka"
	"personalNotificationService/kafka/mocks"
	"personalNotificationService/repositories"
	dbmocks "personalNotificationService/repositories/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter(handler *Handlers) *gin.Engine {
	router := gin.Default()
	router.POST("/send-notification", handler.SendNotification)
	router.POST("/template", handler.CreateTemplate)
	router.POST("/notification", handler.CreateNotification)
	router.POST("/user-notification-channels", handler.CreateUserNotificationChannels)
	return router
}

func TestSendNotification(t *testing.T) {
	mockUserDB := dbmocks.NewMockUserDatabase(t)
	mockNotificationDB := dbmocks.NewMockNotificationDatabase(t)
	mockPriorityProducer := mocks.NewMockPriorityProducer(t)

	handler := &Handlers{
		UserDB:               mockUserDB,
		NotificationDB:       mockNotificationDB,
		NotificationProducer: mockPriorityProducer,
	}

	router := setupRouter(handler)

	username := "johndoe"
	notificationID := uuid.New()
	requestBody := SendNotificationRequest{
		Username:       username,
		NotificationID: notificationID.String(),
	}

	notification := &repositories.Notification{
		ID:       notificationID,
		Type:     repositories.NotificationType("successful_transaction"),
		Priority: repositories.NotificationPriority("high"),
	}

	notificationMessage := &kafka.NotificationMessage{
		Username:       username,
		NotificationID: notificationID.String(),
		Priority:       kafka.MessagePriority(notification.Priority),
	}

	mockUserDB.On("GetUserByUsername", username).Return(&repositories.User{}, nil).Once()
	mockNotificationDB.On("GetNotificationByID", notificationID).Return(notification, nil).Once()
	mockPriorityProducer.On("SendMessage", mock.MatchedBy(func(msg kafka.MessageWithPriority) bool {
		msgParsed := msg.(*kafka.NotificationMessage)
		return msgParsed.Username == notificationMessage.Username &&
			msgParsed.NotificationID == notificationMessage.NotificationID &&
			msgParsed.Priority == notificationMessage.Priority
	})).Return(nil).Once()

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/send-notification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestSendNotification_UserNotFound(t *testing.T) {
	mockUserDB := dbmocks.NewMockUserDatabase(t)
	mockNotificationDB := dbmocks.NewMockNotificationDatabase(t)
	mockPriorityProducer := mocks.NewMockPriorityProducer(t)

	handler := &Handlers{
		UserDB:               mockUserDB,
		NotificationDB:       mockNotificationDB,
		NotificationProducer: mockPriorityProducer,
	}

	router := setupRouter(handler)

	username := "unknownuser"
	notificationID := uuid.New()
	requestBody := SendNotificationRequest{
		Username:       username,
		NotificationID: notificationID.String(),
	}

	mockUserDB.On("GetUserByUsername", username).Return(nil, errors.New("test error")).Once()

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/send-notification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestSendNotification_NotificationNotFound(t *testing.T) {
	mockUserDB := dbmocks.NewMockUserDatabase(t)
	mockNotificationDB := dbmocks.NewMockNotificationDatabase(t)
	mockPriorityProducer := mocks.NewMockPriorityProducer(t)

	handler := &Handlers{
		UserDB:               mockUserDB,
		NotificationDB:       mockNotificationDB,
		NotificationProducer: mockPriorityProducer,
	}

	router := setupRouter(handler)

	username := "johndoe"
	notificationID := uuid.New()
	requestBody := SendNotificationRequest{
		Username:       username,
		NotificationID: notificationID.String(),
	}

	mockUserDB.On("GetUserByUsername", username).Return(&repositories.User{}, nil).Once()
	mockNotificationDB.On("GetNotificationByID", notificationID).Return(nil, errors.New("test error")).Once()

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/send-notification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestSendNotification_ProducerError(t *testing.T) {
	mockUserDB := dbmocks.NewMockUserDatabase(t)
	mockNotificationDB := dbmocks.NewMockNotificationDatabase(t)
	mockPriorityProducer := mocks.NewMockPriorityProducer(t)

	handler := &Handlers{
		UserDB:               mockUserDB,
		NotificationDB:       mockNotificationDB,
		NotificationProducer: mockPriorityProducer,
	}

	router := setupRouter(handler)

	username := "johndoe"
	notificationID := uuid.New()
	requestBody := SendNotificationRequest{
		Username:       username,
		NotificationID: notificationID.String(),
	}

	notification := &repositories.Notification{
		ID:       notificationID,
		Type:     repositories.NotificationType("successful_transaction"),
		Priority: repositories.NotificationPriority("high"),
	}

	mockUserDB.On("GetUserByUsername", username).Return(&repositories.User{}, nil).Once()
	mockNotificationDB.On("GetNotificationByID", notificationID).Return(notification, nil).Once()
	mockPriorityProducer.On("SendMessage", mock.Anything).Return(errors.New("test error")).Once()

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/send-notification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestCreateTemplate(t *testing.T) {
	mockTemplateDB := dbmocks.NewMockTemplateDatabase(t)
	mockTemplateDB.On("CreateTemplate", mock.Anything).Return(nil).Once()

	handler := &Handlers{
		TemplateDB: mockTemplateDB,
	}
	router := setupRouter(handler)

	requestBody := CreateTemplateRequest{
		Channel:  "email",
		Template: "Welcome {{name}}!",
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/template", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	var response TemplateResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, requestBody.Channel, response.Channel)
	assert.Equal(t, requestBody.Template, response.Template)
}

func TestCreateNotification(t *testing.T) {
	mockNotificationDB := dbmocks.NewMockNotificationDatabase(t)
	mockNotificationDB.On("CreateNotification", mock.Anything).Return(nil).Once()

	handler := &Handlers{
		NotificationDB: mockNotificationDB,
	}

	router := setupRouter(handler)

	requestBody := NotificationRequest{
		Type:     "successful_transaction",
		Priority: "high",
		ChannelToTemplateID: []ChannelToTemplate{
			{Channel: "email", TemplateID: uuid.New().String()},
		},
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/notification", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	var response NotificationResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, requestBody.Type, response.Type)
	assert.Equal(t, requestBody.Priority, response.Priority)
	assert.ElementsMatch(t, requestBody.ChannelToTemplateID, response.ChannelToTemplateID)
}

func TestCreateUserNotificationChannels(t *testing.T) {
	mockUserNotificationChannelsDB := dbmocks.NewMockUserNotificationChannelsDatabase(t)
	mockUserNotificationChannelsDB.On("CreateUserNotificationChannel", mock.Anything).Return(nil).Once()

	handler := &Handlers{
		UserNotificationChannelsDB: mockUserNotificationChannelsDB,
	}

	router := setupRouter(handler)

	notificationID := uuid.New()
	requestBody := UserNotificationChannelsRequest{
		Username:       "johndoe",
		NotificationID: notificationID.String(),
		Channels:       []string{"email", "sms"},
	}

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/user-notification-channels", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	var response UserNotificationChannelsResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, requestBody.Username, response.Username)
	assert.Equal(t, requestBody.NotificationID, response.NotificationID)
	assert.ElementsMatch(t, requestBody.Channels, response.Channels)
}
