package main

// go run internal\i18n\gen\main.go

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	sep := string(os.PathSeparator)

	// Path to your translation file
	jsonPath := strings.Join([]string{"internal", "i18n", "translations", "en.json"}, sep)

	dataBytes, err := os.ReadFile(jsonPath)
	if err != nil {
		panic(fmt.Errorf("failed reading json: %w", err))
	}

	var data map[string]interface{}
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		panic(fmt.Errorf("invalid json: %w", err))
	}

	var res []string

	res = append(res, "package i18n")
	res = append(res, "")
	res = append(res, "const (")
	res = append(res, "")

	for key, value := range data {
		// Comment
		res = append(res, fmt.Sprintf("\t// %v", value))
		// Constant name (upper case)
		res = append(res, fmt.Sprintf("\t%s = \"%s\"", strings.ToUpper(key), key))
		res = append(res, "")
	}

	// Remove trailing blank line before ")"
	if len(res) > 0 && strings.TrimSpace(res[len(res)-1]) == "" {
		res = res[:len(res)-1]
	}

	res = append(res, ")")
	res = append(res, "")

	// Output file
	outputPath := strings.Join([]string{"internal", "i18n", "i18n.keys.go"}, sep)

	// Ensure folder exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		panic(fmt.Errorf("failed creating directory: %w", err))
	}

	// Write file
	if err := os.WriteFile(outputPath, []byte(strings.Join(res, "\n")), 0644); err != nil {
		panic(fmt.Errorf("failed writing file: %w", err))
	}

	fmt.Println("Generated:", outputPath)
}
