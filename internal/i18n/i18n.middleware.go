package i18n

import (
	"sort"
	"strconv"
	"strings"

	"github.com/example/internal/config/myconstant"
	"github.com/example/util/arrayutil"
	"github.com/labstack/echo/v4"
)

type TranslateHandlerFunc func(args ...string) string

func (m *I18nModule) I18nMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		genericCode := c.Request().Header.Get("Accept-Language")
		genericCode = m.getCodeFromAcceptLanguage(genericCode)
		viewfinderLangCode := c.Request().Header.Get("lang")
		lang := ""
		if arrayutil.Contains(&m.supportLang, viewfinderLangCode) {
			lang = viewfinderLangCode
		} else {
			lang = genericCode
		}
		var translate TranslateHandlerFunc
		switch lang {
		case Eng:
			translate = m.English.translate
		case Fre:
			translate = m.French.translate
		case Vie:
			translate = m.Vetnamese.translate
		default:
			translate = m.Vetnamese.translate
		}

		c.Set(myconstant.CONTEXT_KEY_TRANSLATION, &translate)
		c.Set(myconstant.CONTEXT_KEY_LANGUAGE_CODE, &lang)

		return next(c)
	}
}

func (m *I18nModule) getCodeFromAcceptLanguage(headerValue string) string {
	type Data struct {
		Lang string
		Q    float64
	}

	lang := []Data{}

	raw := strings.ReplaceAll(headerValue, " ", "")
	raw = strings.ToLower(raw)
	datas := strings.Split(raw, ",")

	for _, v := range datas {
		if v == "" {
			continue
		}
		langValue := strings.Split(v, ";")[0]
		qValue := ""
		if strings.Contains(v, ";") {
			qValue = strings.Split(v, ";")[1]
		}
		q, err := strconv.ParseFloat(strings.ReplaceAll(qValue, "q=", ""), 64)

		if err != nil {
			q = 0
		}

		lang = append(lang, Data{
			Lang: strings.Split(langValue, "-")[0],
			Q:    q,
		})

	}

	sort.Slice(lang, func(i, j int) bool {
		return lang[i].Q > lang[j].Q
	})

	if len(lang) > 0 {
		return lang[0].Lang
	}

	return ""
}
