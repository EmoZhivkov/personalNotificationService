package content

import (
	"bytes"
	"personalNotificationService/common"
	"personalNotificationService/repositories"
	texttemplate "text/template"
)

type slackGenerator struct {
	Params
}

func newSlackGenerator(params Params) (generatorConcreteImpl, error) {
	return &slackGenerator{Params: params}, nil
}

func (e *slackGenerator) GenerateContent(metadata common.Metadata, template repositories.Template) (common.Content, error) {
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
