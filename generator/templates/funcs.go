package templates

import (
	"strings"
	"unicode"
)

func Upper(s string) string {
	return strings.ToUpper(s)
}

func UpperFirst(s string) string {
	// TODO refactor it
	if s == "id" {
		return "ID"
	}

	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
