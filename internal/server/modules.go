package server

import (
	"net/http"

	"github.com/example/internal/config"
	"github.com/example/internal/config/myconstant"
	"github.com/example/internal/i18n"
	"github.com/example/internal/modules/auth"
	"github.com/example/internal/modules/user"
	"github.com/example/internal/service/emailservice"
	"github.com/example/internal/service/jwt"
	"github.com/example/internal/service/myredis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func CreateModule(e *echo.Echo, DB *pgxpool.Pool, env *config.Environment, httpClient *http.Client) {
	// Create Services
	EmailService := emailservice.CreateMailTrapService(env, httpClient)
	JwtService := jwt.NewJwtService(env.JWT_SECRET)
	RedisService := myredis.NewRedisService(env, JwtService)
	i18nModule := i18n.CreateI18nModule()

	// load common to context
	e.Use(i18nModule.I18nMiddleware)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(myconstant.CONTEXT_KEY_REDIS, RedisService)
			c.Set(myconstant.CONTEXT_KEY_ENV, env)
			return next(c)
		}
	})

	// Create modules
	userModule := user.CreateModule(DB)
	authModule := auth.CreateModule(DB)

	// Import to Module
	userModule.Import(RedisService, EmailService)
	authModule.Import(RedisService, userModule.Repo)

	// Register Routes
	root := e.Group("/v" + env.VERSION)
	authModule.RegisterRoutes(root)
	userModule.RegisterRoutes(root)
}
