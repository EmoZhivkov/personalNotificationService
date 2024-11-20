package content

import (
	"bytes"
	htmltemplate "html/template"
	"personalNotificationService/common"
	"personalNotificationService/repositories"
)

type emailGenerator struct {
	Params
}

func newEmailGenerator(params Params) (generatorConcreteImpl, error) {
	return &emailGenerator{Params: params}, nil
}

func (e *emailGenerator) GenerateContent(metadata common.Metadata, template repositories.Template) (common.Content, error) {
	t, err := htmltemplate.New("example").Parse(template.Template)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, metadata); err != nil {
		return "", err
	}
	return common.Content(buf.String()), nil
}
