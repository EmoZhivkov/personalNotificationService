package sender

import (
	"testing"

	"personalNotificationService/common"
	"personalNotificationService/common/mocks"
	"personalNotificationService/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/fiorix/go-smpp/smpp"
)

func TestSmsSender_Send(t *testing.T) {
	config := common.Config{}

	tests := []struct {
		name                     string
		userNotificationSettings any
		content                  common.Content
		mockSetup                func(*mocks.MockSmsClient)
		expectError              bool
	}{
		{
			name: "Valid SMS notification settings",
			userNotificationSettings: repositories.SmsNotificationSettings{
				UserNumber: "1234567890",
			},
			content: "Hello, SMS!",
			mockSetup: func(mockSmsClient *mocks.MockSmsClient) {
				mockSmsClient.On("Submit", mock.AnythingOfType("*smpp.ShortMessage")).Return(&smpp.ShortMessage{}, nil)
			},
			expectError: false,
		},
		{
			name:                     "Invalid notification settings",
			userNotificationSettings: "InvalidSettings",
			content:                  "Hello, SMS!",
			mockSetup:                func(mockSmsClient *mocks.MockSmsClient) {},
			expectError:              true,
		},
		// TODO: uncomment when going to prod, so we return error when client fails
		//{
		//	name: "SMS client fails to send message",
		//	userNotificationSettings: repositories.SmsNotificationSettings{
		//		UserNumber: "1234567890",
		//	},
		//	content: "Hello, SMS!",
		//	mockSetup: func(mockSmsClient *mocks.MockSmsClient) {
		//		mockSmsClient.On("Submit", mock.AnythingOfType("*smpp.ShortMessage")).Return(nil, errors.New("SMPP error"))
		//	},
		//	expectError: true,
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockSmsClient := mocks.NewMockSmsClient(t)
			test.mockSetup(mockSmsClient)

			smsSender := &SmsSender{
				Config:    config,
				SmsClient: mockSmsClient,
			}

			err := smsSender.Send(test.userNotificationSettings, test.content)

			if test.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
