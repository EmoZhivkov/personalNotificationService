package content

import (
	"personalNotificationService/common"
	"personalNotificationService/repositories"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEmailGenerator_GenerateContent(t *testing.T) {
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
				"Name": "John Doe",
			},
			template: repositories.Template{
				ID:       uuid.New(),
				Channel:  "email",
				Template: "Hello, {{.Name}}!",
			},
			expected:    "Hello, John Doe!",
			expectError: false,
		},
		{
			name: "Invalid template syntax",
			metadata: common.Metadata{
				"Name": "John Doe",
			},
			template: repositories.Template{
				ID:       uuid.New(),
				Channel:  "email",
				Template: "Hello, {{.Name}",
			},
			expected:    "",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			generator := &emailGenerator{}
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
