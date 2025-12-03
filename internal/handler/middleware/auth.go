package mymiddleware

import (
	"errors"
	"net/http"

	"github.com/example/internal/config/myconstant"
	"github.com/example/internal/handler/response"
	"github.com/example/internal/i18n"
	"github.com/example/internal/service/myredis"
	"github.com/example/util"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := util.GetTokenString(c)
		if token == "" {
			return response.Error(c, http.StatusUnauthorized, errors.New(i18n.KEY_AUTHORIZATION_NO_TOKEN))
		}

		redisService := util.GetContextValue[*myredis.RedisService](c, myconstant.CONTEXT_KEY_REDIS)

		claims, err := redisService.ValidateToken(c.Request().Context(), token)
		if err != nil {
			redisService.RemoveToken(c.Request().Context(), token)
			if err == redis.Nil {
				return response.Error(c, http.StatusUnauthorized, errors.New(i18n.KEY_AUTHORIZATION_INVALID_CREDENTIAL))
			}
			return response.Error(c, http.StatusUnauthorized, err)
		}

		c.Set(myconstant.CONTEXT_KEY_USER_TYPE, claims.UserType)
		c.Set(myconstant.CONTEXT_KEY_USER_ID, claims.UserID)

		return next(c)
	}
}

func ReadTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := util.GetTokenString(c)
		if token == "" {
			return next(c)
		}

		redisService := util.GetContextValue[*myredis.RedisService](c, myconstant.CONTEXT_KEY_REDIS)

		claims, err := redisService.ValidateToken(c.Request().Context(), token)
		if err != nil {
			redisService.RemoveToken(c.Request().Context(), token)
			return next(c)
		}

		c.Set(myconstant.CONTEXT_KEY_USER_TYPE, claims.UserType)
		c.Set(myconstant.CONTEXT_KEY_USER_ID, claims.UserID)

		return next(c)
	}
}
