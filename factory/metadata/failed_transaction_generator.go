package metadata

import (
	"personalNotificationService/common"
)

type failedTransactionMetadataGenerator struct {
	Params
}

func newFailedTransactionMetadataGenerator(params Params) (generatorConcreteImpl, error) {
	return &failedTransactionMetadataGenerator{Params: params}, nil
}

func (s *failedTransactionMetadataGenerator) GenerateMetadata() (common.Metadata, error) {
	return nil, nil
}
