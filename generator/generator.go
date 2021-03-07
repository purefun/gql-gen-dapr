package generator

import (
	// "bytes"
	"errors"
	"go/format"

	"github.com/purefun/gql-gen-dapr/generator/templates"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

var skipTypes = map[string]bool{
	"__Directive":         true,
	"__DirectiveLocation": true,
	"__Type":              true,
	"__TypeKind":          true,
	"__Field":             true,
	"__EnumValue":         true,
	"__InputValue":        true,
	"__Schema":            true,
}

type Models struct {
	Objects []*Object
}

type Object struct {
	Description string
	Name        string
	Fields      []*Field
	Implements  []string
}

type Field struct {
	Description string
	Name        string
	Type        string
	Tag         string
}

type Generator struct {
	PackageName string
	Sources     []*ast.Source
	Schema      *ast.Schema
	Models      *Models
	ServiceName string
	Imports     []string
}

func (g *Generator) Generate() (string, error) {

	err := g.LoadSchema()
	if err != nil {
		return "", err
	}

	packageOut, err := g.genPackage()
	if err != nil {
		return "", err
	}

	modelsOut, err := g.genModels()
	if err != nil {
		return "", err
	}
	serviceOut, err := g.genService()
	if err != nil {
		return "", err
	}

	importsOut, err := g.genImports()
	if err != nil {
		return "", err
	}

	return packageOut + importsOut + modelsOut + serviceOut, nil
}

func NewSource(name, schemaString string) *ast.Source {
	return &ast.Source{Name: name, Input: schemaString}
}

func (g *Generator) AddSource(name, content string) {
	g.Sources = append(g.Sources, NewSource(name, content))
}

// func (g *Generator) P(ss ...string) {
// 	if len(ss) == 0 {
// 		g.Out.WriteString("\n")
// 		return
// 	}
// 	for _, s := range ss {
// 		g.Out.WriteString(s)
// 	}
// }

func (g *Generator) LoadSchema() error {
	if len(g.Sources) == 0 {
		return errors.New("generator: empty source")
	}
	schema, err := gqlparser.LoadSchema(g.Sources...)
	if err != nil {
		return err
	}
	g.Schema = schema
	return nil
}

func (g *Generator) genPackage() (string, error) {
	out, err := templates.Golang.Execute("package.tmpl", g)
	if err != nil {
		return "", err
	}
	return out, nil
}

func (g *Generator) genModels() (string, error) {
	g.Models = &Models{}

	for _, schemaType := range g.Schema.Types {

		if _, ok := skipTypes[schemaType.Name]; ok {
			continue
		}

		if schemaType == g.Schema.Query ||
			schemaType == g.Schema.Mutation ||
			schemaType == g.Schema.Subscription {
			continue
		}

		switch schemaType.Kind {
		case ast.Object:
			obj := &Object{Name: schemaType.Name}
			for _, field := range schemaType.Fields {
				fieldDefinition := g.Schema.Types[field.Type.Name()]
				switch fieldDefinition.Kind {
				case ast.Scalar:
					obj.Fields = append(obj.Fields, &Field{Name: field.Name, Type: "string"})
				}
			}
			g.Models.Objects = append(g.Models.Objects, obj)

		}
	}

	out, err := templates.Golang.Execute("models.tmpl", g.Models)

	formatted, err := format.Source([]byte(out))

	if err != nil {
		return "", err
	}

	return string(formatted), nil
}

func (g *Generator) genService() (string, error) {
	// g.addImport("context")
	// g.addImport("encoding/json")
	// g.addImport("github.com/dapr/go-sdk/client")
	// g.addImport("github.com/dapr/go-sdk/service/common")
	return "", nil
}

func (g *Generator) addImport(s string) {
	g.Imports = append(g.Imports, s)
}

func (g *Generator) genImports() (string, error) {
	out, err := templates.Golang.Execute("imports.tmpl", g)
	if err != nil {
		return "", err
	}
	return out, nil
}
