package user

import (
	mymiddleware "github.com/example/internal/handler/middleware"
	"github.com/example/internal/service/emailservice"
	"github.com/example/internal/service/myredis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type UserModule struct {
	Handler *UserHandler
	Repo    *UserRepository
}

func CreateModule(db *pgxpool.Pool) *UserModule {
	repo := UserRepository{DB: db}
	controller := UserHandler{Repo: &repo}
	return &UserModule{Repo: &repo, Handler: &controller}
}

func (m *UserModule) Import(
	RedisService myredis.RedisServiceInterface,
	EmailService emailservice.EmailServiceInterface,
) {
	m.Handler.RedisService = RedisService
	m.Handler.EmailService = EmailService
}

func (m *UserModule) RegisterRoutes(group *echo.Group) {
	user := group.Group("/user")
	user.POST("/reset-password", m.Handler.ResetPassword)
	user.POST("/request-password-reset-email", m.Handler.RequestPasswordResetEmail)
	user.POST("/request-email-verification", m.Handler.RequestEmailVerification)
	user.POST("/verify-email", m.Handler.VerifyEmail)

	userAuth := user.Group("", mymiddleware.AuthMiddleware)
	userAuth.GET("/detail", m.Handler.detail)
	userAuth.POST("/update", m.Handler.Update)
	userAuth.POST("/update-password", m.Handler.UpdatePassword)
	userAuth.POST("/update-username", m.Handler.UpdateUserName)
}
