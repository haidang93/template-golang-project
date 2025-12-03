package user

import (
	"github.com/example/internal/config/myconstant"
	"github.com/example/internal/i18n"
	"github.com/example/internal/service/emailservice"
	"github.com/example/util"
	"github.com/labstack/echo/v4"
)

func CreateEmailVerificationTemplate(c echo.Context, user *User, code *string) *emailservice.SendData {
	LangCode := util.GetContextValue[*string](c, myconstant.CONTEXT_KEY_LANGUAGE_CODE)
	templateID := ""
	switch *LangCode {
	case i18n.Eng:
		templateID = "3c8eb226-0aab-44ce-8c12-63f3b9ecc895"
	case i18n.Fre:
		templateID = "0cc39b6c-25df-48dd-8ecc-be4047d89284"
	default:
		templateID = "25c4be57-3bb5-45da-a091-5d9ab9838d19"
	}

	fromEmail := "no-reply@domain.vn"
	fromName := "Example"
	return &emailservice.SendData{
		TemplateUUID: &templateID,
		From: &emailservice.Email{
			Email: &fromEmail,
			Name:  &fromName,
		},
		To: &[]emailservice.Email{
			{
				Email: user.Email,
				Name:  user.DisplayName,
			},
		},
		TemplateVariables: &map[string]interface{}{
			"name": user.FirstName,
			"code": code,
		},
	}
}

func CreateResetPasswordTemplate(c echo.Context, user *User, code *string) *emailservice.SendData {
	LangCode := util.GetContextValue[*string](c, myconstant.CONTEXT_KEY_LANGUAGE_CODE)
	templateID := ""
	switch *LangCode {
	case i18n.Eng:
		templateID = "b0ec7a10-2899-41b1-8fc5-09cd92c8ed94"
	case i18n.Fre:
		templateID = "afa060d1-e403-4f25-9ec0-daa9e9065da5"
	default:
		templateID = "229223ad-a6b4-4167-b535-d1db12ca567c"
	}

	fromEmail := "no-reply@domain.vn"
	fromName := "Example"
	return &emailservice.SendData{
		TemplateUUID: &templateID,
		From: &emailservice.Email{
			Email: &fromEmail,
			Name:  &fromName,
		},
		To: &[]emailservice.Email{
			{
				Email: user.Email,
				Name:  user.DisplayName,
			},
		},
		TemplateVariables: &map[string]interface{}{
			"name": user.FirstName,
			"code": code,
		},
	}
}
