package templates

import (
	"bytes"
	"text/template"

	"github.com/GeertJohan/go.rice"
)

type TemplateBox struct {
	box *rice.Box
}

func (t *TemplateBox) Execute(tplName string, data interface{}) (string, error) {
	tplStr, err := t.box.String(tplName)
	if err != nil {
		return "", err
	}
	tpl, err := template.New(tplName).Parse(tplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	err = tpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

var Golang = &TemplateBox{rice.MustFindBox("golang")}
