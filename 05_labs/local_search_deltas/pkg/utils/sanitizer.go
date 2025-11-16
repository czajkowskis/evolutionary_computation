package utils

import "strings"

// SanitizeFileName normalises a string so it can be safely used as a filename.
func SanitizeFileName(name string) string {
	replacer := strings.NewReplacer(" ", "_", "(", "", ")", "", ",", "")
	return replacer.Replace(name)
}
