package repositories

import "github.com/google/uuid"

type TemplateDatabase interface {
	CreateTemplate(template *Template) error
	CreateTemplates(templates Templates) error
	GetTemplateByID(id uuid.UUID) (*Template, error)
	GetTemplatesByIDs(ids uuid.UUIDs) (Templates, error)
	GetTemplatesByChannel(channel NotificationChannel) ([]Template, error)
}

type templateDatabase struct {
	client DbClient
}

func NewTemplateDatabase(client DbClient) TemplateDatabase {
	return &templateDatabase{client: client}
}

func (db *templateDatabase) CreateTemplate(template *Template) error {
	if err := db.client.Create(template).Error; err != nil {
		return err
	}
	return nil
}

func (db *templateDatabase) CreateTemplates(templates Templates) error {
	if err := db.client.Create(&templates).Error; err != nil {
		return err
	}
	return nil
}

func (db *templateDatabase) GetTemplateByID(id uuid.UUID) (*Template, error) {
	var template Template
	if err := db.client.Take(&template, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (p *templateDatabase) GetTemplatesByIDs(ids uuid.UUIDs) (Templates, error) {
	var templates Templates
	if err := p.client.Find(&templates, "id IN ?", ids).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (db *templateDatabase) GetTemplatesByChannel(channel NotificationChannel) ([]Template, error) {
	var templates []Template
	if err := db.client.Find(&templates, "channel = ?", channel).Error; err != nil {
		return nil, err
	}
	return templates, nil
}
