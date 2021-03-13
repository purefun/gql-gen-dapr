package generator

import (
	"github.com/purefun/gql-gen-dapr/tools"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestData struct {
	name string
	g    *Generator
	file string
}

func g() *Generator {
	return NewGenerator(Options{
		PackageName: "testdata",
		ServiceName: "Example",
	})
}

var tests = []TestData{
	{
		name: "package name",
		file: "./testdata/package-name.gql",
		g:    g(),
	},
	{
		name: "scalar-fields",
		file: "./testdata/scalar-fields.gql",
		g:    g(),
	},
	{
		name: "resolvers",
		file: "./testdata/resolvers.gql",
		g:    g(),
	},
	{
		name: "descriptions",
		file: "./testdata/descriptions.gql",
		g:    g(),
	},
	{
		name: "reference fields",
		file: "./testdata/ref-fields.gql",
		g:    g(),
	},
	{
		name: "enums",
		file: "./testdata/enums.gql",
		g:    g(),
	},
	{
		name: "unions",
		file: "./testdata/unions.gql",
		g:    g(),
	},
	{
		name: "interfces",
		file: "./testdata/interfaces.gql",
		g:    g(),
	},
	{
		name: "directive tag",
		file: "./testdata/directive-tag.gql",
		g:    g(),
	},
}

func TestGenerate(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schemaString, err := ioutil.ReadFile(tt.file)
			require.NoError(t, err)

			tt.g.AddSource(tt.file, string(schemaString))

			goFile := tools.ReplaceExt(tt.file, ".go")
			goBytes, err := ioutil.ReadFile(goFile)
			require.NoError(t, err)

			out, err := tt.g.Generate()
			require.NoError(t, err)

			require.Equal(t, strings.TrimSpace(string(goBytes)), strings.TrimSpace(out))
		})
	}
}
