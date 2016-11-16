package utils

import (
	"regexp"
	"strings"
)

// RemoveHashtags removes all '#' characters
func RemoveHashtags(text string) string {
	// TODO: don't suppress errors
	r, _ := regexp.Compile("#([^\\s]*)")
	return strings.TrimSpace(r.ReplaceAllString(text, ""))
}

// FindHashtags finds all tags starting with '#'
func FindHashtags(text string) []string {
	// TODO: confirm this and don't suppress errors
	r, _ := regexp.Compile(`/(\s|^)\#\w\w+\b/gm`)
	return r.FindAllString(text, -1)
}
