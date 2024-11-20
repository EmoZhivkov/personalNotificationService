package content

import (
	"bytes"
	"personalNotificationService/common"
	"personalNotificationService/repositories"
	texttemplate "text/template"
)

type smsGenerator struct {
	Params
}

func newSmsGenerator(params Params) (generatorConcreteImpl, error) {
	return &smsGenerator{Params: params}, nil
}

func (e *smsGenerator) GenerateContent(metadata common.Metadata, template repositories.Template) (common.Content, error) {
	t, err := texttemplate.New("example").Parse(template.Template)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, metadata); err != nil {
		return "", err
	}
	return common.Content(buf.String()), nil
}
