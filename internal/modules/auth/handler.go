package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/example/internal/handler/response"
	"github.com/example/internal/i18n"
	"github.com/example/internal/models"
	"github.com/example/internal/modules/user"
	"github.com/example/internal/pkg/validate"
	"github.com/example/internal/service/myredis"
	"github.com/example/util"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type AuthHandler struct {
	Repo         *AuthRepository
	redisService myredis.RedisServiceInterface
	UserService  user.UserRepositoryInterface
}

func (ctr *AuthHandler) Signup(c echo.Context) error {
	var req SignupDto
	if err := validate.Bind(c, &req); err != nil {
		return response.BadRequestErr(c, err)
	}

	if !IsValidEmail(req.Email) {
		return response.BadRequest(c, i18n.KEY_VALIDATE_NOT_EMAIL)
	}

	if err := user.PasswordStrengthCheck(req.Password); err != nil {
		return response.BadRequestErr(c, err)
	}

	hashedPassword, err := user.HashPassword(req.Password)
	if err != nil {
		return response.BadRequestErr(c, err)
	}

	foundUser, err := ctr.Repo.checkEmailExist(req.Email)
	if err != nil {
		return response.BadRequestErr(c, err)
	}
	if *foundUser {
		return response.BadRequest(c, i18n.KEY_DUPLICATE_CONSTRAINT_EMAIL)
	}

	if req.Username == nil {
		req.Username = ParseUsernameFromEmail(req.Email)
	}

	err = user.ValidateUsername(req.Username)
	if err != nil {
		return response.BadRequestErr(c, err)
	}

	email := strings.ToLower(*req.Email)
	totp, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "ViewFinder",
		AccountName: email,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return response.InternalServerErrorErr(c, err)
	}
	secret := totp.Secret()

	userType := user.USER_TYPE_NORMAL
	defaultStatus := models.COMMON_STATUS_OPEN
	dto := user.User{
		Password:     &hashedPassword,
		Username:     req.Username,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        &email,
		UserType:     &userType,
		Status:       &defaultStatus,
		ReceiveEmail: req.ReceiveEmail,
		AcceptPolicy: req.AcceptPolicy,
		Secret:       &secret,
	}
	newUser, err := ctr.UserService.Create(&dto)
	if err != nil {
		return response.BadRequestErr(c, err)
	}

	return c.JSON(http.StatusOK, newUser)
}

func (ctr *AuthHandler) Signin(c echo.Context) error {
	var req SigninDto
	if err := validate.Bind(c, &req); err != nil {
		return response.BadRequestErr(c, err)
	}

	ArrayEmail := []string{req.Email}
	dto := user.UserRequestDto{
		ArrayEmail: &ArrayEmail,
	}
	foundUser, err := ctr.UserService.GetOne(&dto)
	if err != nil {
		if err == pgx.ErrNoRows {
			return response.Unauthorized(c, i18n.KEY_AUTHORIZATION_WRONG_EMAIL_OR_PASSWORD)
		}
		return response.BadRequestErr(c, err)
	}

	checkPass := user.CheckPasswordHash(&req.Password, foundUser.Password)
	if !checkPass {
		return response.Unauthorized(c, i18n.KEY_AUTHORIZATION_WRONG_EMAIL_OR_PASSWORD)
	}

	if foundUser.EmailVerified == nil || !*foundUser.EmailVerified {
		return response.Unauthorized(c, "You account's email is not verified")
	}

	userType := user.USER_TYPE_NORMAL
	token, err := ctr.redisService.CreateToken(c.Request().Context(), foundUser.ID, foundUser.Email, &userType)
	if err != nil {
		return response.BadRequestErr(c, err)
	}

	currentTime := time.Now()
	updateData := user.User{
		LastLogin: &currentTime,
	}
	updated, err := ctr.UserService.Update(foundUser.ID, &updateData)
	if err != nil {
		return response.BadRequestErr(c, err)
	}

	return response.Success(c, map[string]any{"accessToken": token, "user": updated})
}

func (ctr *AuthHandler) RefreshToken(c echo.Context) error {
	UserID := util.GetUserID(c)
	dto := user.UserRequestDto{
		ArrayID: &[]string{*UserID},
	}
	foundUser, err := ctr.UserService.GetOne(&dto)
	if err != nil {
		if err == pgx.ErrNoRows {
			return response.NotFound(c, i18n.KEY_NOT_FOUND_USER)
		}
		return response.BadRequestErr(c, err)
	}

	if *foundUser.Status != models.COMMON_STATUS_OPEN {
		return response.BadRequest(c, i18n.KEY_BAD_REQUEST_USER_NOT_ACTIVE)
	}

	userType := user.USER_TYPE_NORMAL
	token, err := ctr.redisService.CreateToken(c.Request().Context(), foundUser.ID, foundUser.Email, &userType)
	if err != nil {
		return response.BadRequestErr(c, err)
	}

	if err := ctr.redisService.RemoveToken(c.Request().Context(), util.GetTokenString(c)); err != nil {
		return response.BadRequestErr(c, err)
	}

	return c.JSON(http.StatusOK, map[string]any{"accessToken": token})
}
