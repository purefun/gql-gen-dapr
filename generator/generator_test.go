package generator

import (
	"io/ioutil"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestData struct {
	name string
	g    Generator
	file string
}

func newGenerator() Generator {
	return Generator{PackageName: "testdata", ServiceName: "Example"}
}

var tests = []TestData{
	{
		name: "package name",
		file: "./testdata/package-name.gql",
		g:    newGenerator(),
	},
	{
		name: "scalar-fields",
		file: "./testdata/scalar-fields.gql",
		g:    newGenerator(),
	},
	{
		name: "resolvers",
		file: "./testdata/resolvers.gql",
		g:    newGenerator(),
	},
	{
		name: "descriptions",
		file: "./testdata/descriptions.gql",
		g:    newGenerator(),
	},
}

func TestGenerate(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schemaString, err := ioutil.ReadFile(tt.file)
			require.NoError(t, err)

			tt.g.AddSource(tt.file, string(schemaString))

			var re = regexp.MustCompile(`^(.+)(\.gql)`)
			goFile := re.ReplaceAllString(tt.file, `$1.go`)
			goBytes, err := ioutil.ReadFile(goFile)
			require.NoError(t, err)

			out, err := tt.g.Generate()
			require.NoError(t, err)

			require.Equal(t, strings.TrimSpace(string(goBytes)), strings.TrimSpace(out))
		})
	}
}
