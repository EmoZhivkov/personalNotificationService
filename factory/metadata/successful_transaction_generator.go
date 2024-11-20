package metadata

import (
	"personalNotificationService/common"
)

type successfulTransactionMetadataGenerator struct {
	Params
}

func newSuccessfulTransactionMetadataGenerator(params Params) (generatorConcreteImpl, error) {
	return &successfulTransactionMetadataGenerator{Params: params}, nil
}

func (s *successfulTransactionMetadataGenerator) GenerateMetadata() (common.Metadata, error) {
	return nil, nil
}
