package api

import (
	"log"
	"personalNotificationService/auth"
	"personalNotificationService/common"
	"personalNotificationService/kafka"
	"personalNotificationService/repositories"

	_ "personalNotificationService/docs" // docs is generated by Swag CLI

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

type Server struct {
	hostPort string
	router   *gin.Engine
}

func NewServer(config common.Config, notificationsProducer kafka.PriorityProducer) *Server {
	router := gin.Default()
	router.Use(recoverFromPanicMiddleware)

	client, err := repositories.NewDbClient(config.DbConnectionURL)
	if err != nil {
		log.Fatalf("Error creating Database client: %v", err)
	}

	userDatabase := repositories.NewUserDatabase(client)
	notificationDatabase := repositories.NewNotificationDatabase(client)
	templateDatabase := repositories.NewTemplateDatabase(client)
	userNotificationChannelsDatabase := repositories.NewUserNotificationChannelsDatabase(client)

	handlers := NewHandlers(notificationsProducer, userDatabase, templateDatabase, notificationDatabase, userNotificationChannelsDatabase)

	secured := router.Group("/api/v1")
	secured.Use(validateJWTMiddleware(config.JwtSecret))
	{
		secured.POST("/notifications/send", handlers.SendNotification)
		secured.POST("/templates", handlers.CreateTemplate)
		secured.POST("/notifications", handlers.CreateNotification)
		secured.POST("/user-notification-channels", handlers.CreateUserNotificationChannels)
	}

	// TODO: remove when going to prod
	setupExampleData(userDatabase, notificationDatabase, templateDatabase, userNotificationChannelsDatabase)

	authHandlers := auth.NewAuthenticationHandler(config.JwtSecret, userDatabase)
	router.POST("/api/v1/authenticate", authHandlers.Authenticate)

	// http://localhost:8080/swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &Server{
		hostPort: config.HostPort,
		router:   router,
	}
}

func (s *Server) Run() error {
	log.Printf("Starting Notification API on port: %s", s.hostPort)
	return s.router.Run(":" + s.hostPort)
}

func createTestNotifications(notificationDatabase repositories.NotificationDatabase, testTemplates repositories.Templates) repositories.Notifications {
	testNotifications := repositories.Notifications{
		{
			ID:       uuid.New(),
			Type:     repositories.SuccessfulTransactionNotificationType,
			Priority: repositories.HighNotificationPriority,
			ChannelToTemplateID: map[repositories.NotificationChannel]uuid.UUID{
				repositories.EmailNotificationChannel: testTemplates[0].ID,
				repositories.SlackNotificationChannel: testTemplates[1].ID,
				repositories.SmsNotificationChannel:   testTemplates[2].ID,
			},
		},
		{
			ID:       uuid.New(),
			Type:     repositories.FailedTransactionNotificationType,
			Priority: repositories.MediumNotificationPriority,
			ChannelToTemplateID: map[repositories.NotificationChannel]uuid.UUID{
				repositories.EmailNotificationChannel: testTemplates[0].ID,
				repositories.SlackNotificationChannel: testTemplates[1].ID,
				repositories.SmsNotificationChannel:   testTemplates[2].ID,
			},
		},
	}

	// Check and insert notifications
	if err := notificationDatabase.CreateNotifications(testNotifications); err != nil {
		log.Fatalf("Error creating test notifications: %v", err)
	}

	notifications, err := notificationDatabase.GetNotificationsByIDs(uuid.UUIDs{testNotifications[0].ID, testNotifications[1].ID})
	if err != nil {
		log.Fatalf("Error getting test notifications: %v", err)
	}

	log.Printf("%v", notifications)

	return testNotifications
}

func createTestTemplates(templateDatabase repositories.TemplateDatabase) repositories.Templates {
	testTemplates := repositories.Templates{
		{
			ID:       uuid.New(),
			Channel:  repositories.EmailNotificationChannel,
			Template: "Hello, this is an email notification template.",
		},
		{
			ID:       uuid.New(),
			Channel:  repositories.SlackNotificationChannel,
			Template: "Hello, this is a Slack notification template.",
		},
		{
			ID:       uuid.New(),
			Channel:  repositories.SmsNotificationChannel,
			Template: "Hello, this is an SMS notification template.",
		},
	}

	// Check and insert templates
	if err := templateDatabase.CreateTemplates(testTemplates); err != nil {
		log.Fatalf("Error creating test templates: %v", err)
	}

	templates, err := templateDatabase.GetTemplatesByIDs(uuid.UUIDs{testTemplates[0].ID, testTemplates[1].ID, testTemplates[2].ID})
	if err != nil {
		log.Fatalf("Error getting test templates: %v", err)
	}

	log.Printf("%v", templates)
	return testTemplates
}

func createTestUserNotificationChannels(userNotificationChannelsDatabase repositories.UserNotificationChannelsDatabase, testUsers repositories.Users, testNotifications repositories.Notifications) {
	var testChannels []repositories.UserNotificationChannels

	for _, user := range testUsers {
		for _, notification := range testNotifications {
			testChannels = append(testChannels, repositories.UserNotificationChannels{
				Username:       user.Username,
				NotificationID: notification.ID,
				Channels: []repositories.NotificationChannel{
					repositories.EmailNotificationChannel,
					repositories.SlackNotificationChannel,
					repositories.SmsNotificationChannel,
				},
			})
		}
	}

	// Bulk insert channels
	if err := userNotificationChannelsDatabase.BulkCreateUserNotificationChannel(testChannels); err != nil {
		log.Fatalf("Error creating test user notification channels: %v", err)
	}

	for _, user := range testUsers {
		for _, notification := range testNotifications {
			notificationChannels, err := userNotificationChannelsDatabase.GetUserNotificationChannels(user.Username, notification.ID)
			if err != nil {
				log.Fatalf("Error getting user notification channels: %v", err)
			}
			log.Printf("%v", notificationChannels)
		}
	}
}

func createTestUsers(userDatabase repositories.UserDatabase) repositories.Users {
	testUsers := repositories.Users{
		{Username: "alice", PasswordHash: "alice", NotificationSettings: map[repositories.NotificationChannel]interface{}{
			repositories.EmailNotificationChannel: repositories.EmailNotificationSettings{UserEmail: "alice"},
			repositories.SlackNotificationChannel: repositories.SlackNotificationSettings{UserHandle: "alice@alice.com"},
			repositories.SmsNotificationChannel:   repositories.SmsNotificationSettings{UserNumber: "123132123"},
		}},
		{Username: "bob", PasswordHash: "bob", NotificationSettings: map[repositories.NotificationChannel]interface{}{
			repositories.EmailNotificationChannel: repositories.EmailNotificationSettings{UserEmail: "bob"},
			repositories.SlackNotificationChannel: repositories.SlackNotificationSettings{UserHandle: "bob@bob.com"},
			repositories.SmsNotificationChannel:   repositories.SmsNotificationSettings{UserNumber: "312132312"},
		}},
		{Username: "carol", PasswordHash: "carol", NotificationSettings: map[repositories.NotificationChannel]interface{}{
			repositories.EmailNotificationChannel: repositories.EmailNotificationSettings{UserEmail: "carol"},
			repositories.SlackNotificationChannel: repositories.SlackNotificationSettings{UserHandle: "carol@carol.com"},
			repositories.SmsNotificationChannel:   repositories.SmsNotificationSettings{UserNumber: "312123123"},
		}},
		{Username: "dave", PasswordHash: "dave", NotificationSettings: map[repositories.NotificationChannel]interface{}{
			repositories.EmailNotificationChannel: repositories.EmailNotificationSettings{UserEmail: "dave"},
			repositories.SlackNotificationChannel: repositories.SlackNotificationSettings{UserHandle: "dave@dave.com"},
			repositories.SmsNotificationChannel:   repositories.SmsNotificationSettings{UserNumber: "31223123"},
		}},
	}

	testUsernames := make([]string, 0, len(testUsers))
	for i := 0; i < len(testUsers); i++ {
		hashedPassword, err := auth.HashPassword(testUsers[i].PasswordHash)
		if err != nil {
			log.Fatalf("Error hashing password: %v", err)
		}
		testUsers[i].PasswordHash = hashedPassword

		testUsernames = append(testUsernames, testUsers[i].Username)
	}

	existingTestUsers, err := userDatabase.GetUsersByUsernames(testUsernames)
	if err != nil {
		log.Fatalf("Error getting existing test users: %v", err)
	}
	existingTestUsersSet := existingTestUsers.ToUserSet()

	usersToCreate := repositories.Users{}
	for _, testUser := range testUsers {
		if _, exists := existingTestUsersSet[testUser.Username]; !exists {
			usersToCreate = append(usersToCreate, testUser)
		}
	}

	if len(usersToCreate) > 0 {
		if err := userDatabase.CreateUsers(usersToCreate); err != nil {
			log.Fatalf("Error creating test users: %v", err)
		}
	}

	return testUsers
}

func setupExampleData(
	userDatabase repositories.UserDatabase,
	notificationDatabase repositories.NotificationDatabase,
	templateDatabase repositories.TemplateDatabase,
	userNotificationChannelsDatabase repositories.UserNotificationChannelsDatabase,
) {
	users := createTestUsers(userDatabase)

	testTemplates := createTestTemplates(templateDatabase)
	testNotifications := createTestNotifications(notificationDatabase, testTemplates)
	createTestUserNotificationChannels(userNotificationChannelsDatabase, users, testNotifications)
}
