package db

import (
	"regexp"
	"strings"
)

func toSnakeCase(str string) string {
	pattern := regexp.MustCompile(`([A-Z])`)
	result := pattern.ReplaceAllStringFunc(str, func(match string) string {
		return "_" + strings.ToLower(match)
	})
	return strings.TrimPrefix(result, "_")
}
