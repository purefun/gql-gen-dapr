package generator

import (
	"bytes"

	"github.com/purefun/gql-gen-dapr/generator/templates"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

type Generator struct {
	PackageName string
	Sources     []*ast.Source
	Out         *bytes.Buffer
}

func (g *Generator) Generate() (string, error) {
	g.Out = &bytes.Buffer{}
	schema, err := gqlparser.LoadSchema(g.Sources...)
	if err != nil {
		return "", err
	}

	for _, t := range schema.Types {
		switch t.Kind {
		case ast.Object:

		}
	}

	out, tplErr := templates.Golang.Execute("object.tmpl", struct{}{})
	if tplErr != nil {
		return "", tplErr
	}

	g.WriteOut(out)

	return g.Out.String(), nil
}

func NewSource(name, schemaString string) *ast.Source {
	return &ast.Source{Name: name, Input: schemaString}
}

func (g *Generator) AddSource(name, content string) {
	g.Sources = append(g.Sources, NewSource(name, content))
}

func (g *Generator) WriteOut(o string) {
	g.Out.WriteString(o)
}
