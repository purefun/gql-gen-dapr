package generator

import (
	"errors"
	"fmt"
	"go/format"
	"sort"
	"strings"
	"unicode"

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

	"__schema": true,
	"__type":   true,
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
	NonNull     bool
	Tag         string
}

type Enum struct {
	Description string
	Name        string
	Values      []*EnumValue
}

type EnumValue struct {
	Description string
	Name        string
}

var ScalarMap = map[string]string{
	"ID":      "string",
	"String":  "string",
	"Int":     "int",
	"Float":   "float64",
	"Boolean": "bool",
}

type Resolver struct {
	Name string
	Type string
}

type Query struct {
	Resolvers []*Resolver
}

type Generator struct {
	PackageName string
	Sources     []*ast.Source
	Schema      *ast.Schema
	Models      *Models
	ServiceName string
	Imports     map[string]string // package->alias
	Query       *Query
	Enums       []*Enum
}

func (g *Generator) Generate() (string, error) {
	if g.Imports == nil {
		g.Imports = make(map[string]string)
	}

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

	out := packageOut + importsOut + modelsOut + serviceOut

	formatted, err := format.Source([]byte(out))

	if err != nil {
		return "", fmt.Errorf("format source failed: %w, source: %s", err, out)
	}

	return string(formatted), nil
}

func NewSource(name, schemaString string) *ast.Source {
	return &ast.Source{Name: name, Input: schemaString}
}

func (g *Generator) AddSource(name, content string) {
	g.Sources = append(g.Sources, NewSource(name, content))
}

func (g *Generator) LoadSchema() error {
	if len(g.Sources) == 0 {
		return errors.New("no source")
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

	sortedKeys := []string{}
	for key := range g.Schema.Types {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	for _, typeKey := range sortedKeys {
		schemaType, ok := g.Schema.Types[typeKey]
		if !ok {
			continue
		}

		if _, ok := skipTypes[schemaType.Name]; ok {
			continue
		}

		if schemaType == g.Schema.Query ||
			schemaType == g.Schema.Mutation ||
			schemaType == g.Schema.Subscription {
			continue
		}

		switch schemaType.Kind {
		case ast.Object, ast.InputObject:
			obj := &Object{Name: schemaType.Name, Description: schemaType.Description}
			for _, field := range schemaType.Fields {
				fieldDefinition := g.Schema.Types[field.Type.Name()]

				switch fieldDefinition.Kind {
				case ast.Scalar:
					scalarName, ok := ScalarMap[field.Type.NamedType]
					if !ok {
						panic(fmt.Errorf("invalid graphql scalar name: %s", field.Type.NamedType))
					}
					obj.Fields = append(obj.Fields, &Field{
						Name:        field.Name,
						Type:        scalarName,
						NonNull:     field.Type.NonNull,
						Description: field.Description,
					})
				case ast.Object, ast.Enum:
					obj.Fields = append(obj.Fields, &Field{
						Name:        field.Name,
						Type:        field.Type.Name(),
						NonNull:     field.Type.NonNull,
						Description: field.Description,
					})
				}
			}
			g.Models.Objects = append(g.Models.Objects, obj)

		case ast.Enum:
			g.addImport("fmt", "")
			g.addImport("io", "")
			g.addImport("strconv", "")

			e := &Enum{
				Name:        schemaType.Name,
				Description: schemaType.Description,
			}
			for _, v := range schemaType.EnumValues {
				e.Values = append(e.Values, &EnumValue{
					Name:        v.Name,
					Description: v.Description,
				})
			}

			g.Enums = append(g.Enums, e)
		}
	}

	out, err := templates.Golang.Execute("models.tmpl", g)
	if err != nil {
		return "", err
	}

	return out, nil
}

func (g *Generator) genService() (string, error) {
	g.Query = &Query{Resolvers: []*Resolver{}}

	if g.Schema.Query != nil {
		g.addImport("context", "")
		g.addImport("encoding/json", "")
		g.addImport("github.com/dapr/go-sdk/client", "")
		g.addImport("github.com/dapr/go-sdk/service/common", "")

		for _, field := range g.Schema.Query.Fields {

			if _, ok := skipTypes[field.Name]; ok {
				continue
			}
			r := &Resolver{Name: field.Name, Type: strings.ToLower(field.Type.Name())}
			g.Query.Resolvers = append(g.Query.Resolvers, r)
		}

		out, err := templates.Golang.Execute("service.tmpl", g)
		if err != nil {
			return "", err
		}
		return out, nil
	}
	return "", nil
}

func (g *Generator) addImport(pkg, alias string) {
	g.Imports[pkg] = alias
}

func (g *Generator) genImports() (string, error) {
	out, err := templates.Golang.Execute("imports.tmpl", g)
	if err != nil {
		return "", err
	}
	return out, nil
}

func (g *Generator) GoName(name string) string {
	if name == "id" {
		return "ID"
	}
	r := []rune(name)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func (g *Generator) GoDoc(name, desc string) string {
	if desc == "" {
		return ""
	}
	n := g.GoName(name)
	return "// " + n + " " + strings.Replace(desc, "\n", "\n"+"// ", -1)
}
