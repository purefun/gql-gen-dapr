package generator

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestData struct {
	name    string
	file    string
	options Options
}

func TestGenerate(t *testing.T) {
	tests := []TestData{
		{
			name: "package name",
			file: "package-name",
			options: Options{
				PkgName: "main",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gqlString, err := ioutil.ReadFile("./testdata/" + tt.file + ".gql")
			assert.NoError(t, err)

			goString, err := ioutil.ReadFile("./testdata/" + tt.file + ".go")
			assert.NoError(t, err)

			out := Generate(string(gqlString), tt.options)
			want := string(goString)

			assert.Equal(t, want, out)
		})
	}
}
