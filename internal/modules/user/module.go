package user

import (
	mymiddleware "github.com/example/internal/handler/middleware"
	"github.com/example/internal/service/emailservice"
	"github.com/example/internal/service/myredis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type UserModule struct {
	Controller *UserController
	Repo       *UserRepository
}

func CreateModule(db *pgxpool.Pool) *UserModule {
	repo := UserRepository{DB: db}
	controller := UserController{Repo: &repo}
	return &UserModule{Repo: &repo, Controller: &controller}
}

func (m *UserModule) Import(
	RedisService myredis.RedisServiceInterface,
	EmailService emailservice.EmailServiceInterface,
) {
	m.Controller.RedisService = RedisService
	m.Controller.EmailService = EmailService
}

func (m *UserModule) RegisterRoutes(group *echo.Group) {
	user := group.Group("/user")
	user.POST("/reset-password", m.Controller.ResetPassword)
	user.POST("/request-password-reset-email", m.Controller.RequestPasswordResetEmail)
	user.POST("/request-email-verification", m.Controller.RequestEmailVerification)
	user.POST("/verify-email", m.Controller.VerifyEmail)

	userAuth := user.Group("", mymiddleware.AuthMiddleware)
	userAuth.GET("/detail", m.Controller.detail)
	userAuth.POST("/update", m.Controller.Update)
	userAuth.POST("/update-password", m.Controller.UpdatePassword)
	userAuth.POST("/update-username", m.Controller.UpdateUserName)
}
