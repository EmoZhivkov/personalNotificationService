package content

import (
	"testing"

	"personalNotificationService/common"
	"personalNotificationService/repositories"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSlackGenerator_GenerateContent(t *testing.T) {
	tests := []struct {
		name        string
		metadata    common.Metadata
		template    repositories.Template
		expected    string
		expectError bool
	}{
		{
			name: "Valid template and metadata",
			metadata: common.Metadata{
				"Name": "Jane Doe",
			},
			template: repositories.Template{
				ID:       uuid.New(),
				Channel:  "slack",
				Template: "Hello, {{.Name}}!",
			},
			expected:    "Hello, Jane Doe!",
			expectError: false,
		},
		{
			name: "Invalid template syntax",
			metadata: common.Metadata{
				"Name": "Jane Doe",
			},
			template: repositories.Template{
				ID:       uuid.New(),
				Channel:  "slack",
				Template: "Hello, {{.Name",
			},
			expected:    "",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			generator := &slackGenerator{}
			content, err := generator.GenerateContent(test.metadata, test.template)

			if test.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.expected, string(content))
		})
	}
}
