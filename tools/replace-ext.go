package tools

import (
	"fmt"
	"regexp"
	"strings"
)

func ReplaceExt(f, newExt string) string {
	re := regexp.MustCompile(`(.+)\.[a-zA-Z0-9]+$`)
	ext := strings.TrimPrefix(newExt, ".")
	return re.ReplaceAllString(f, fmt.Sprintf("$1.%s", ext))
}
