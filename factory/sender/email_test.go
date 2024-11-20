package sender

import (
	"personalNotificationService/common/mocks"
	"testing"

	"personalNotificationService/common"
	"personalNotificationService/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEmailSender_Send(t *testing.T) {
	config := common.Config{}

	tests := []struct {
		name                     string
		userNotificationSettings any
		content                  common.Content
		mockSetup                func(*mocks.MockMailClient)
		expectError              bool
	}{
		{
			name: "Valid email notification settings",
			userNotificationSettings: repositories.EmailNotificationSettings{
				UserEmail: "test@example.com",
			},
			content: "Hello, World!",
			mockSetup: func(mockMailClient *mocks.MockMailClient) {
				mockMailClient.On("DialAndSend", mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name:                     "Invalid notification settings",
			userNotificationSettings: "InvalidSettings",
			content:                  "Hello, World!",
			mockSetup:                func(mockMailClient *mocks.MockMailClient) {},
			expectError:              true,
		},
		// TODO: uncomment when going to prod, so we return error when client fails
		//{
		//	name: "Mail client fails to send email",
		//	userNotificationSettings: repositories.EmailNotificationSettings{
		//		UserEmail: "test@example.com",
		//	},
		//	content: "Hello, World!",
		//	mockSetup: func(mockMailClient *mocks.MockMailClient) {
		//		mockMailClient.On("DialAndSend", mock.Anything).Return(errors.New("SMTP error"))
		//	},
		//	expectError: true,
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockMailClient := mocks.NewMockMailClient(t)
			test.mockSetup(mockMailClient)

			emailSender := &EmailSender{
				Config:     config,
				MailClient: mockMailClient,
			}

			err := emailSender.Send(test.userNotificationSettings, test.content)

			if test.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
