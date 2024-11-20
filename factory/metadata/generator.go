package metadata

import (
	"errors"
	"personalNotificationService/common"
	"personalNotificationService/repositories"
)

type generatorConcreteImpl interface {
	GenerateMetadata() (common.Metadata, error)
}

var metadataGeneratorFactoryMethods = map[repositories.NotificationType]func(params Params) (generatorConcreteImpl, error){
	repositories.SuccessfulTransactionNotificationType: newSuccessfulTransactionMetadataGenerator,
	repositories.FailedTransactionNotificationType:     newFailedTransactionMetadataGenerator,
}

type Generator interface {
	GenerateMetadata(repositories.NotificationType) (common.Metadata, error)
}

func NewMetadataGenerator(params Params) Generator {
	return &generator{Params: params}
}

type Params struct{}

type generator struct {
	Params
}

func (g *generator) GenerateMetadata(notificationType repositories.NotificationType) (common.Metadata, error) {
	factoryMethod, ok := metadataGeneratorFactoryMethods[notificationType]
	if !ok {
		return nil, errors.New("error invalid notification type")
	}

	metadataGenerator, err := factoryMethod(g.Params)
	if err != nil {
		return nil, errors.New("error instantiating sender concrete implementation")
	}

	return metadataGenerator.GenerateMetadata()
}
