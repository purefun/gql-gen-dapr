package generator

import (
	"errors"
	"fmt"
	"go/format"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/purefun/gql-gen-dapr/generator/templates"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

const Version = "v0.4.6"

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

type Interface struct {
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

type Argument struct {
	Name string
	Type string
}

type Resolver struct {
	Name     string
	Type     string
	Argument *Argument
}

type Query struct {
	Resolvers []*Resolver
}

type Options struct {
	PackageName       string
	ServiceName       string
	GenHeaderComments bool
}

func NewGenerator(o Options) *Generator {
	return &Generator{
		PackageName:       o.PackageName,
		ServiceName:       o.ServiceName,
		Version:           Version,
		GenHeaderComments: o.GenHeaderComments,
	}
}

type Generator struct {
	Version           string
	PackageName       string
	GenHeaderComments bool
	Sources           []*ast.Source
	Schema            *ast.Schema
	Models            *Models
	ServiceName       string
	// TODO move to models
	Imports    map[string]string // package->alias
	Query      *Query
	Enums      []*Enum
	Interfaces []*Interface
}

func (g *Generator) Generate() (string, error) {
	if g.Imports == nil {
		g.Imports = make(map[string]string)
	}

	g.addTagDirective()

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
		return "", fmt.Errorf("format source failed: %w, source: \n%s", err, out)
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
					// tag directive
					tag := ""
					for _, directive := range field.Directives {
						if directive.Name == "tag" {
							for _, arg := range directive.Arguments {
								tag += arg.Name + ":" + "\"" + arg.Value.Raw + "\" "
							}
						}
					}
					obj.Fields = append(obj.Fields, &Field{
						Name:        field.Name,
						Type:        scalarName,
						NonNull:     field.Type.NonNull,
						Description: field.Description,
						Tag:         strings.TrimSpace(tag),
					})
				case ast.Object, ast.Enum:
					typeName := field.Type.String()
					isArray, _ := regexp.MatchString(`^\[.+\]!?$`, typeName)
					if isArray {
						typeName = "[]" + field.Type.Name()
					}
					obj.Fields = append(obj.Fields, &Field{
						Name:        field.Name,
						Type:        typeName,
						NonNull:     field.Type.NonNull,
						Description: field.Description,
					})
				}
			}
			for _, impl := range g.Schema.GetImplements(schemaType) {
				obj.Implements = append(obj.Implements, impl.Name)
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

		case ast.Union, ast.Interface:
			it := &Interface{
				Name:        schemaType.Name,
				Description: schemaType.Description,
			}
			g.Interfaces = append(g.Interfaces, it)
		}

	}

	out, err := templates.Golang.Execute("models.tmpl", g)
	if err != nil {
		return "", err
	}

	return out, nil
}

func (g *Generator) genService() (string, error) {
	if g.Schema.Query == nil && g.Schema.Mutation == nil {
		return "", nil
	}

	g.addImport("context", "")
	g.addImport("encoding/json", "")
	g.addImport("github.com/dapr/go-sdk/client", "")
	g.addImport("github.com/dapr/go-sdk/service/common", "")
	g.addImport("github.com/dapr/go-sdk/service/grpc", "")

	fields := ast.FieldList{}

	if g.Schema.Query != nil {
		fields = append(fields, g.Schema.Query.Fields...)
	}
	if g.Schema.Mutation != nil {
		fields = append(fields, g.Schema.Mutation.Fields...)
	}

	g.Query = &Query{Resolvers: []*Resolver{}}

	for _, field := range fields {
		if _, ok := skipTypes[field.Name]; ok {
			continue
		}
		typeName, ok := ScalarMap[field.Type.NamedType]
		if !ok {
			typeName = field.Type.NamedType
		}

		r := &Resolver{Name: field.Name, Type: typeName}

		if len(field.Arguments) > 1 {
			panic("the number of resolver arguments should be 0 or 1")
		}

		if len(field.Arguments) == 1 {
			arg := field.Arguments[0]
			typeName, ok := ScalarMap[arg.Type.NamedType]
			if !ok {
				typeName = arg.Type.NamedType
			}
			r.Argument = &Argument{Name: arg.Name, Type: typeName}
		}

		g.Query.Resolvers = append(g.Query.Resolvers, r)
	}

	out, err := templates.Golang.Execute("service.tmpl", g)
	if err != nil {
		return "", err
	}
	return out, nil
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

func (g *Generator) addTagDirective() {
	tagSource := &ast.Source{
		Name:    "tag.graphql",
		Input:   "directive @tag(any: String) on FIELD_DEFINITION",
		BuiltIn: true,
	}

	g.Sources = append(g.Sources, tagSource)
}
