package templates

import (
	"bytes"
	"text/template"

	"github.com/GeertJohan/go.rice"
)

type TemplateBox struct {
	box *rice.Box
}

type Option func(funcs template.FuncMap)

func (t *TemplateBox) Execute(tplName string, data interface{}, opts ...Option) (string, error) {

	for _, o := range opts {
		o(funcs)
	}

	tplStr, err := t.box.String(tplName)
	if err != nil {
		return "", err
	}
	tpl, err := template.New(tplName).Funcs(funcs).Parse(tplStr)
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
