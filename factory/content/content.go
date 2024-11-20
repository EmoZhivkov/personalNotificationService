package content

import (
	"errors"
	"fmt"
	"personalNotificationService/common"
	"personalNotificationService/repositories"
)

type generatorConcreteImpl interface {
	GenerateContent(common.Metadata, repositories.Template) (common.Content, error)
}

var contentGeneratorFactoryMethods = map[repositories.NotificationChannel]func(Params) (generatorConcreteImpl, error){
	repositories.EmailNotificationChannel: newEmailGenerator,
	repositories.SlackNotificationChannel: newSlackGenerator,
	repositories.SmsNotificationChannel:   newSmsGenerator,
}

type Generator interface {
	GenerateContent(common.Metadata, repositories.Template) (common.Content, error)
}

func NewContentGenerator(params Params) Generator {
	return &generator{Params: params}
}

type Params struct{}

type generator struct {
	Params
}

func (g *generator) GenerateContent(metadata common.Metadata, template repositories.Template) (common.Content, error) {
	factoryMethod, ok := contentGeneratorFactoryMethods[template.Channel]
	if !ok {
		return "", fmt.Errorf("error unsupported notification channel: %v", template.Channel)
	}

	contentGenerator, err := factoryMethod(g.Params)
	if err != nil {
		return "", errors.New("error instantiating content generator concrete implementation")
	}

	return contentGenerator.GenerateContent(metadata, template)
}
