package i18n

import (
	"regexp"
	"strings"
)

func Translate(translations map[string]string, args ...string) string {
	key := args[0]
	args = args[1:]
	translate := translations[key]

	// Predefined key not found. Use raw value instead
	if translate == "" && key != "" {
		translate = key
	}

	placeHolders := extractUniquePlaceholders(translate)

	for i, value := range placeHolders {
		translate = strings.ReplaceAll(translate, "{"+value+"}", args[i])
	}

	return translate
}

func extractUniquePlaceholders(s string) []string {
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(s, -1)

	seen := make(map[string]bool)
	var result []string

	for _, m := range matches {
		placeholder := m[1]
		if !seen[placeholder] {
			seen[placeholder] = true
			result = append(result, placeholder)
		}
	}
	return result
}
