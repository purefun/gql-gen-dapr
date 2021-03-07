package generator

import (
	"go/format"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestData struct {
	name      string
	g         Generator
	file      string
	assertion assert.ErrorAssertionFunc
}

func newGen() Generator {
	return Generator{PackageName: "testdata", ServiceName: "Example"}
}

var tests = []TestData{
	{
		name:      "package name",
		file:      "./testdata/package-name.gql",
		g:         newGen(),
		assertion: assert.NoError,
	},
	{
		name:      "scalar-fields",
		file:      "./testdata/scalar-fields.gql",
		g:         newGen(),
		assertion: assert.NoError,
	},
	{
		name:      "basic service",
		file:      "./testdata/resolvers.gql",
		g:         newGen(),
		assertion: assert.NoError,
	},
}

func TestGenerate(t *testing.T) {

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schemaString, err := ioutil.ReadFile(tt.file)
			assert.NoError(t, err)

			tt.g.AddSource(tt.file, string(schemaString))

			var re = regexp.MustCompile(`^(.+)(\.gql)`)
			goFile := re.ReplaceAllString(tt.file, `$1.go`)
			goBytes, err := ioutil.ReadFile(goFile)
			assert.NoError(t, err)

			out, err := tt.g.Generate()

			want, err := format.Source(goBytes)
			assert.NoError(t, err)

			if tt.assertion != nil {
				tt.assertion(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, strings.TrimSpace(string(want)), strings.TrimSpace(out))
		})
	}
}
