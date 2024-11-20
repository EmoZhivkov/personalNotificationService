package sender

import (
	"testing"

	"personalNotificationService/common"
	"personalNotificationService/common/mocks"
	"personalNotificationService/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSlackSender_Send(t *testing.T) {
	config := common.Config{
		SlackChannelID: "test-channel",
	}

	tests := []struct {
		name                     string
		userNotificationSettings any
		content                  common.Content
		mockSetup                func(*mocks.MockSlackClient)
		expectError              bool
	}{
		{
			name: "Valid Slack notification settings",
			userNotificationSettings: repositories.SlackNotificationSettings{
				UserHandle: "test-channel",
			},
			content: "Hello, Slack!",
			mockSetup: func(mockSlackClient *mocks.MockSlackClient) {
				mockSlackClient.On("PostMessage", config.SlackChannelID, mock.Anything).Return("timestamp", "channel", nil)
			},
			expectError: false,
		},
		{
			name:                     "Invalid notification settings",
			userNotificationSettings: "InvalidSettings",
			content:                  "Hello, Slack!",
			mockSetup:                func(mockSlackClient *mocks.MockSlackClient) {},
			expectError:              true,
		},
		// TODO: uncomment when going to prod, so we return error when client fails
		//{
		//	name: "Slack client fails to send message",
		//	userNotificationSettings: repositories.SlackNotificationSettings{
		//		UserHandle: "test-channel",
		//	},
		//	content: "Hello, Slack!",
		//	mockSetup: func(mockSlackClient *mocks.MockSlackClient) {
		//		mockSlackClient.On("PostMessage", config.SlackChannelID, mock.Anything).Return("", "", errors.New("Slack API error"))
		//	},
		//	expectError: true,
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockSlackClient := mocks.NewMockSlackClient(t)
			test.mockSetup(mockSlackClient)

			slackSender := &SlackSender{
				Config:      config,
				SlackClient: mockSlackClient,
			}

			err := slackSender.Send(test.userNotificationSettings, test.content)

			if test.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
