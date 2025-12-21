package utils

import "strings"

func SanitizeFileName(name string) string {
	replacer := strings.NewReplacer(" ", "_", "(", "", ")", "", ",", "")
	return replacer.Replace(name)
}

