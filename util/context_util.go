package util

import (
	"strings"

	"github.com/example/internal/config/myconstant"
	"github.com/labstack/echo/v4"
)

// key is [CONTEXT_KEY]
func GetContextValue[T any](c echo.Context, key string) T {
	value, _ := c.Get(key).(T)
	return value
}

func GetTokenString(c echo.Context) string {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		token = c.QueryParam("token")
	}
	return strings.ReplaceAll(token, "Bearer ", "")
}

func GetUserID(c echo.Context) *string {
	return GetContextValue[*string](c, myconstant.CONTEXT_KEY_USER_ID)
}

func GetUserType(c echo.Context) *string {
	return GetContextValue[*string](c, myconstant.CONTEXT_KEY_USER_TYPE)
}

func GetCommunityID(c echo.Context) *string {
	ID := c.Param("communityID")
	if ID == "" {
		return nil
	}
	return &ID
}
