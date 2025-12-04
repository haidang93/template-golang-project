package auth

import (
	mymiddleware "github.com/example/internal/handler/middleware"
	"github.com/example/internal/modules/user"
	"github.com/example/internal/service/myredis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type AuthModule struct {
	handler *AuthHandler
	Repo    *AuthRepository
}

func CreateModule(db *pgxpool.Pool) *AuthModule {
	repo := AuthRepository{DB: db}
	controller := AuthHandler{Repo: &repo}
	return &AuthModule{Repo: &repo, handler: &controller}
}

func (m *AuthModule) Import(
	redisService myredis.RedisServiceInterface,
	UserService user.UserRepositoryInterface,
) {
	m.handler.UserService = UserService
	m.handler.redisService = redisService
}

func (m *AuthModule) RegisterRoutes(group *echo.Group) {
	auth := group.Group("/auth")
	auth.POST("/signin", m.handler.Signin)
	auth.POST("/signup", m.handler.Signup)

	privateAuth := auth.Group("", mymiddleware.AuthMiddleware)
	privateAuth.POST("/refresh-token", m.handler.RefreshToken)
}
