package repositories

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateTemplate(t *testing.T) {
	testTemplateDB := NewTemplateDatabase(getTestInMemoryDBClient())

	template := &Template{
		ID:       uuid.New(),
		Channel:  EmailNotificationChannel,
		Template: "This is a test email template",
	}

	assert.NoError(t, testTemplateDB.CreateTemplate(template))

	retrievedTemplate, err := testTemplateDB.GetTemplateByID(template.ID)
	assert.NoError(t, err)
	assert.Equal(t, template.ID, retrievedTemplate.ID)
	assert.Equal(t, template.Channel, retrievedTemplate.Channel)
	assert.Equal(t, template.Template, retrievedTemplate.Template)
}

func TestCreateTemplates(t *testing.T) {
	testTemplateDB := NewTemplateDatabase(getTestInMemoryDBClient())

	templates := Templates{
		{
			ID:       uuid.New(),
			Channel:  SmsNotificationChannel,
			Template: "This is a test SMS template",
		},
		{
			ID:       uuid.New(),
			Channel:  SlackNotificationChannel,
			Template: "This is a test Slack template",
		},
	}
	assert.NoError(t, testTemplateDB.CreateTemplates(templates))

	template1, err := testTemplateDB.GetTemplateByID(templates[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, templates[0].ID, template1.ID)
	assert.Equal(t, templates[0].Channel, template1.Channel)
	assert.Equal(t, templates[0].Template, template1.Template)

	template2, err := testTemplateDB.GetTemplateByID(templates[1].ID)
	assert.NoError(t, err)
	assert.Equal(t, templates[1].ID, template2.ID)
	assert.Equal(t, templates[1].Channel, template2.Channel)
	assert.Equal(t, templates[1].Template, template2.Template)
}

func TestGetTemplateByID(t *testing.T) {
	testTemplateDB := NewTemplateDatabase(getTestInMemoryDBClient())

	nonExistentID := uuid.New()
	_, err := testTemplateDB.GetTemplateByID(nonExistentID)
	assert.Error(t, err)
}

func TestGetTemplatesByIDs(t *testing.T) {
	testTemplateDB := NewTemplateDatabase(getTestInMemoryDBClient())

	templates := Templates{
		{
			ID:       uuid.New(),
			Channel:  EmailNotificationChannel,
			Template: "Email template 1",
		},
		{
			ID:       uuid.New(),
			Channel:  EmailNotificationChannel,
			Template: "Email template 2",
		},
	}
	assert.NoError(t, testTemplateDB.CreateTemplates(templates))

	retrievedTemplates, err := testTemplateDB.GetTemplatesByIDs([]uuid.UUID{templates[0].ID, templates[1].ID})
	assert.NoError(t, err)
	assert.Len(t, retrievedTemplates, len(templates))

	retrievedTemplateSet := map[uuid.UUID]Template{}
	for _, retrievedTemplate := range retrievedTemplates {
		tmp := retrievedTemplate
		retrievedTemplateSet[retrievedTemplate.ID] = tmp
	}

	template1 := retrievedTemplateSet[templates[0].ID]
	assert.Equal(t, templates[0].ID, template1.ID)
	assert.Equal(t, templates[0].Channel, template1.Channel)
	assert.Equal(t, templates[0].Template, template1.Template)

	template2 := retrievedTemplateSet[templates[1].ID]
	assert.Equal(t, templates[1].ID, template2.ID)
	assert.Equal(t, templates[1].Channel, template2.Channel)
	assert.Equal(t, templates[1].Template, template2.Template)
}

func TestGetTemplatesByChannel(t *testing.T) {
	testTemplateDB := NewTemplateDatabase(getTestInMemoryDBClient())

	templates := Templates{
		{
			ID:       uuid.New(),
			Channel:  EmailNotificationChannel,
			Template: "Email template 1",
		},
		{
			ID:       uuid.New(),
			Channel:  SmsNotificationChannel,
			Template: "SMS template 1",
		},
		{
			ID:       uuid.New(),
			Channel:  EmailNotificationChannel,
			Template: "Email template 2",
		},
	}

	assert.NoError(t, testTemplateDB.CreateTemplates(templates))

	// Retrieve templates by channel
	emailTemplates, err := testTemplateDB.GetTemplatesByChannel(EmailNotificationChannel)
	assert.NoError(t, err)
	assert.Len(t, emailTemplates, 2)

	for _, template := range emailTemplates {
		assert.Equal(t, EmailNotificationChannel, template.Channel)
	}

	smsTemplates, err := testTemplateDB.GetTemplatesByChannel(SmsNotificationChannel)
	assert.NoError(t, err)
	assert.Len(t, smsTemplates, 1)
	assert.Equal(t, SmsNotificationChannel, smsTemplates[0].Channel)
}
