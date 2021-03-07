package generator

import (
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

var tests = []TestData{
	{
		name:      "package name",
		file:      "./testdata/package-name.gql",
		g:         Generator{PackageName: "testdata"},
		assertion: assert.NoError,
	},
	{
		name:      "basic types",
		file:      "./testdata/basic-types.gql",
		g:         Generator{PackageName: "testdata"},
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
			goString, err := ioutil.ReadFile(goFile)
			assert.NoError(t, err)

			out, err := tt.g.Generate()

			want := string(goString)

			if tt.assertion != nil {
				tt.assertion(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, strings.TrimSpace(want), strings.TrimSpace(out))

		})
	}
}
