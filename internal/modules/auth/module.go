package auth

import (
	mymiddleware "github.com/example/internal/handler/middleware"
	"github.com/example/internal/modules/user"
	"github.com/example/internal/service/myredis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type AuthModule struct {
	Controller *AuthController
	Repo       *AuthRepository
}

func CreateModule(db *pgxpool.Pool) *AuthModule {
	repo := AuthRepository{DB: db}
	controller := AuthController{Repo: &repo}
	return &AuthModule{Repo: &repo, Controller: &controller}
}

func (m *AuthModule) Import(
	redisService myredis.RedisServiceInterface,
	UserService user.UserRepositoryInterface,
) {
	m.Controller.UserService = UserService
	m.Controller.redisService = redisService
}

func (m *AuthModule) RegisterRoutes(group *echo.Group) {
	auth := group.Group("/auth")
	auth.POST("/signin", m.Controller.Signin)
	auth.POST("/signup", m.Controller.Signup)

	privateAuth := auth.Group("", mymiddleware.AuthMiddleware)
	privateAuth.POST("/refresh-token", m.Controller.RefreshToken)
}
