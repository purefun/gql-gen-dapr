package templates

import (
	"strings"
	"text/template"
	"unicode"
)

var funcs = template.FuncMap{
	"upper":       upper,
	"upperFirst":  upperFirst,
	"prefixLines": prefixLines,
}

func upper(s string) string {
	return strings.ToUpper(s)
}

func upperFirst(s string) string {
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func prefixLines(prefix, s string) string {
	if s == "" {
		return ""
	}
	return prefix + strings.Replace(s, "\n", "\n"+prefix, -1)
}
