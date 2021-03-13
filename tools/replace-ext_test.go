package tools

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReplaceExt(t *testing.T) {
	type args struct {
		f      string
		newExt string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "replaces extension w/o dot",
			args: args{f: "schema.graphql", newExt: "go"},
			want: "schema.go",
		},
		{
			name: "replaces extension w/ dot",
			args: args{f: "schema.graphql", newExt: ".go"},
			want: "schema.go",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, ReplaceExt(tt.args.f, tt.args.newExt))
		})
	}
}
