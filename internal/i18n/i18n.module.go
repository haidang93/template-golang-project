package i18n

import (
	_ "embed"

	"github.com/example/internal/config/myconstant"
	"github.com/labstack/echo/v4"
)

// support languages
const (
	Vie string = "vi" // fallback language
	Eng string = "en"
	Fre string = "fr"
)

//go:embed translations/vi.json
var viJSON []byte

//go:embed translations/en.json
var enJSON []byte

//go:embed translations/fr.json
var frJSON []byte

type I18nModule struct {
	supportLang []string
	Vetnamese   *translation
	English     *translation
	French      *translation
}

func CreateI18nModule() *I18nModule {
	return &I18nModule{
		Vetnamese:   createTranslation(viJSON),
		English:     createTranslation(enJSON),
		French:      createTranslation(frJSON),
		supportLang: []string{Vie, Eng, Fre},
	}
}

func Trsl(c echo.Context, args ...string) string {
	translate := *c.Get(myconstant.CONTEXT_KEY_TRANSLATION).(*TranslateHandlerFunc)
	return translate(args...)
}
