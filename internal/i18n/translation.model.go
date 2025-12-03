package i18n

import (
	_ "embed"
	"encoding/json"
)

type translation struct {
	translations map[string]string
}

func (m *translation) translate(args ...string) string {
	return Translate(m.translations, args...)
}

func createTranslation(data []byte) *translation {
	var viTranslations map[string]string
	json.Unmarshal(data, &viTranslations)
	return &translation{translations: viTranslations}
}
