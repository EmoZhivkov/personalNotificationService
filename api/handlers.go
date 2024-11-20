package api

import (
	"net/http"
	"personalNotificationService/kafka"

	"personalNotificationService/common"
	"personalNotificationService/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handlers struct {
	NotificationProducer       kafka.PriorityProducer
	UserDB                     repositories.UserDatabase
	TemplateDB                 repositories.TemplateDatabase
	NotificationDB             repositories.NotificationDatabase
	UserNotificationChannelsDB repositories.UserNotificationChannelsDatabase
}

func NewHandlers(
	notificationProducer kafka.PriorityProducer,
	userDB repositories.UserDatabase,
	templateDB repositories.TemplateDatabase,
	notificationDB repositories.NotificationDatabase,
	userNotificationChannelsDB repositories.UserNotificationChannelsDatabase,
) Handlers {
	return Handlers{
		NotificationProducer:       notificationProducer,
		UserDB:                     userDB,
		TemplateDB:                 templateDB,
		NotificationDB:             notificationDB,
		UserNotificationChannelsDB: userNotificationChannelsDB,
	}
}

// SendNotification godoc
// @Summary Send a notification to a user
// @Description Sends a notification to a user based on the notification ID and username
// @Tags Notifications
// @Accept  json
// @Produce  json
// @Param SendNotificationRequest body SendNotificationRequest true "Notification details"
// @Success 200
// @Failure 400 {object} common.ErrorResponse "Invalid Request"
// @Failure 500 {object} common.ErrorResponse "Internal Server Error"
// @Router /api/v1/notifications/send [post]
func (h *Handlers) SendNotification(c *gin.Context) {
	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if _, err := h.UserDB.GetUserByUsername(req.Username); err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	notification, err := h.NotificationDB.GetNotificationByID(uuid.MustParse(req.NotificationID))
	if err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	notificationMessage := req.ToKafkaNotificationMessage(notification)
	if err := h.NotificationProducer.SendMessage(notificationMessage); err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to send notification message")
		return
	}

	c.Status(http.StatusOK)
}

// CreateTemplate godoc
// @Summary Create a notification template
// @Description Creates a template for a specific notification channel
// @Tags Templates
// @Accept  json
// @Produce  json
// @Param CreateTemplateRequest body CreateTemplateRequest true "Template details"
// @Success 201 {object} TemplateResponse
// @Failure 400 {object} common.ErrorResponse "Invalid Request"
// @Failure 500 {object} common.ErrorResponse "Internal Server Error"
// @Router /api/v1/templates [post]
func (h *Handlers) CreateTemplate(c *gin.Context) {
	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	template := req.ToTemplateModel()
	if err := h.TemplateDB.CreateTemplate(&template); err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to create template")
		return
	}

	c.JSON(http.StatusCreated, ToTemplateResponse(&template))
}

// CreateNotification godoc
// @Summary Create a notification
// @Description Creates a notification and maps it to channels and templates
// @Tags Notifications
// @Accept  json
// @Produce  json
// @Param NotificationRequest body NotificationRequest true "Notification details"
// @Success 201 {object} NotificationResponse
// @Failure 400 {object} common.ErrorResponse "Invalid Request"
// @Failure 500 {object} common.ErrorResponse "Internal Server Error"
// @Router /api/v1/notifications [post]
func (h *Handlers) CreateNotification(c *gin.Context) {
	var req NotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	notification := req.ToNotificationModel()
	if err := h.NotificationDB.CreateNotification(&notification); err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to create notification")
		return
	}

	c.JSON(http.StatusCreated, ToNotificationResponse(&notification))
}

// CreateUserNotificationChannels godoc
// @Summary Configure user notification channels
// @Description Configures notification channels for a user
// @Tags Notifications
// @Accept  json
// @Produce  json
// @Param UserNotificationChannelsRequest body UserNotificationChannelsRequest true "User notification channel details"
// @Success 201 {object} UserNotificationChannelsResponse
// @Failure 400 {object} common.ErrorResponse "Invalid Request"
// @Failure 500 {object} common.ErrorResponse "Internal Server Error"
// @Router /api/v1/user-notification-channels [post]
func (h *Handlers) CreateUserNotificationChannels(c *gin.Context) {
	var req UserNotificationChannelsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		common.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	userNotificationChannels := req.ToUserNotificationChannelsModel()
	if err := h.UserNotificationChannelsDB.CreateUserNotificationChannel(&userNotificationChannels); err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to create user notification channels")
		return
	}

	c.JSON(http.StatusCreated, ToUserNotificationChannelsResponse(&userNotificationChannels))
}
