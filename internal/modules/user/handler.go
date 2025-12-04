package user

import (
	"net/http"
	"strconv"
	"time"

	"github.com/example/internal/config/myconstant"
	"github.com/example/internal/handler/response"
	"github.com/example/internal/i18n"
	"github.com/example/internal/pkg/validate"
	"github.com/example/internal/service/emailservice"
	"github.com/example/internal/service/myredis"
	"github.com/example/util"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type UserHandler struct {
	Repo         *UserRepository
	RedisService myredis.RedisServiceInterface
	EmailService emailservice.EmailServiceInterface
}

func (ctr *UserHandler) detail(c echo.Context) error {
	UserID := util.GetUserID(c)

	user, err := ctr.Repo.GetOne(&UserRequestDto{ArrayID: &[]string{*UserID}})
	if err != nil {
		if err == pgx.ErrNoRows {
			return response.BadRequest(c, i18n.KEY_NOT_FOUND_USER)
		}
		return response.BadRequestErr(c, err)
	}

	return c.JSON(http.StatusOK, user)
}

func (ctr *UserHandler) Update(c echo.Context) error {
	UserID := util.GetUserID(c)

	var dto UserUpdateDto
	err := validate.Bind(c, &dto)
	if err != nil {
		return response.BadRequestErr(c, err)
	}

	currentTime := time.Now()
	data := User{
		DisplayName:       dto.DisplayName,
		FirstName:         dto.FirstName,
		LastName:          dto.LastName,
		PhoneNumber:       dto.PhoneNumber,
		Bio:               dto.Bio,
		PreferredLanguage: dto.PreferredLanguage,
		ReceiveEmail:      dto.ReceiveEmail,
		UpdateDate:        &currentTime,
	}

	updated, err := ctr.Repo.Update(UserID, &data)
	if err != nil {
		return response.BadRequestErr(c, err)
	}

	return c.JSON(http.StatusOK, updated)
}

func (ctr *UserHandler) UpdatePassword(c echo.Context) error {
	UserID := util.GetUserID(c)

	type Req struct {
		OldPassword     *string `json:"oldPassword"`
		NewPassword     *string `json:"newPassword"`
		ResetCredential *bool   `json:"resetCredential"`
	}
	var req Req
	err := validate.Bind(c, &req)
	if err != nil {
		return response.BadRequestErr(c, err)
	}

	if err := PasswordStrengthCheck(req.NewPassword); err != nil {
		return response.BadRequestErr(c, err)
	}

	foundUser, err := ctr.Repo.GetOne(&UserRequestDto{ArrayID: &[]string{*UserID}})
	if err != nil {
		return response.InternalServerErrorErr(c, err)
	}

	checkPass := CheckPasswordHash(req.OldPassword, foundUser.Password)
	if !checkPass {
		return response.Forbidden(c, i18n.KEY_VALIDATE_WRONG_PASSWORD)
	}

	hashedPassword, err := HashPassword(req.NewPassword)
	if err != nil {
		return response.InternalServerErrorErr(c, err)
	}

	currentTime := time.Now()
	data := User{
		Password:   &hashedPassword,
		UpdateDate: &currentTime,
	}

	updated, err := ctr.Repo.Update(UserID, &data)
	if err != nil {
		return response.InternalServerErrorErr(c, err)
	}

	var token *string
	if req.ResetCredential != nil && *req.ResetCredential {
		if err := ctr.RedisService.RemoveAllToken(c.Request().Context(), UserID); err != nil {
			return response.InternalServerErrorErr(c, err)
		}
		userType := USER_TYPE_NORMAL
		newToken, err := ctr.RedisService.CreateToken(c.Request().Context(), updated.ID, updated.Email, &userType)
		if err != nil {
			return response.InternalServerErrorErr(c, err)
		}

		token = newToken
	}

	return c.JSON(http.StatusOK, map[string]any{"accessToken": token, "user": updated})
}

func (ctr *UserHandler) UpdateUserName(c echo.Context) error {
	UserID := util.GetUserID(c)

	type Req struct {
		Username *string `json:"username"`
	}
	var req Req
	if err := validate.Bind(c, &req); err != nil {
		return response.BadRequestErr(c, err)
	}

	err := ValidateUsername(req.Username)
	if err != nil {
		return response.BadRequestErr(c, err)
	}

	currentTime := time.Now()
	data := User{
		Username:   req.Username,
		UpdateDate: &currentTime,
	}

	updated, err := ctr.Repo.Update(UserID, &data)
	if err != nil {
		return response.BadRequestErr(c, err)
	}

	return c.JSON(http.StatusOK, updated)
}

func (ctr *UserHandler) ResetPassword(c echo.Context) error {
	type Req struct {
		Email           *string `json:"email"`
		Code            *string `json:"code"`
		ResetCredential *bool   `json:"resetCredential"`
		NewPassword     *string `json:"newPassword"`
	}
	var req Req
	if err := validate.Bind(c, &req); err != nil {
		return response.BadRequestErr(c, err)
	}

	user, err := ctr.Repo.GetOne(&UserRequestDto{ArrayEmail: &[]string{*req.Email}})
	if err != nil {
		if err == pgx.ErrNoRows {
			return response.BadRequest(c, i18n.KEY_NOT_FOUND_USER)
		}
		return response.BadRequestErr(c, err)
	}

	if user.EmailVerified == nil || !*user.EmailVerified {
		return response.Unauthorized(c, i18n.KEY_AUTHORIZATION_EMAIL_NOT_VERIFIED)
	}

	if user.Secret == nil {
		return response.InternalServerError(c, i18n.KEY_INTERNAL_ERROR_NO_SECRET_CODE)
	}

	ok, err := totp.ValidateCustom(*req.Code, *user.Secret, time.Now(), totp.ValidateOpts{
		Period:    myconstant.PASSCODE_EXPIRATION_DURATION,
		Algorithm: otp.AlgorithmSHA1,
		Digits:    6,
	})
	if err != nil {
		return response.InternalServerErrorErr(c, err)
	}

	if !ok {
		return response.Unauthorized(c, i18n.KEY_INTERNAL_ERROR_WRONG_PASSCODE)
	}

	if err := PasswordStrengthCheck(req.NewPassword); err != nil {
		return response.BadRequestErr(c, err)
	}

	hashedPassword, err := HashPassword(req.NewPassword)
	if err != nil {
		return response.InternalServerErrorErr(c, err)
	}

	dto := User{
		Password: &hashedPassword,
	}
	data, err := ctr.Repo.Update(user.ID, &dto)
	if err != nil {
		return response.InternalServerErrorErr(c, err)
	}

	if req.ResetCredential != nil && *req.ResetCredential {
		if err := ctr.RedisService.RemoveAllToken(c.Request().Context(), user.ID); err != nil {
			return response.InternalServerErrorErr(c, err)
		}
	}

	userType := USER_TYPE_NORMAL
	newToken, err := ctr.RedisService.CreateToken(c.Request().Context(), data.ID, data.Email, &userType)
	if err != nil {
		return response.InternalServerErrorErr(c, err)
	}

	return response.Success(c, map[string]any{"accessToken": newToken, "user": data})
}

func (ctr *UserHandler) VerifyEmail(c echo.Context) error {
	type Req struct {
		Email *string `json:"email"`
		Code  *string `json:"code"`
	}
	var req Req
	if err := validate.Bind(c, &req); err != nil {
		return response.BadRequestErr(c, err)
	}

	user, err := ctr.Repo.GetOne(&UserRequestDto{ArrayEmail: &[]string{*req.Email}})
	if err != nil {
		if err == pgx.ErrNoRows {
			return response.BadRequest(c, i18n.KEY_NOT_FOUND_USER)
		}
		return response.BadRequestErr(c, err)
	}

	if user.EmailVerified != nil && *user.EmailVerified {
		return response.BadRequest(c, i18n.KEY_BAD_REQUEST_ALREADY_VERIFIED)
	}

	if user.Secret == nil {
		return response.InternalServerError(c, i18n.KEY_INTERNAL_ERROR_NO_SECRET_CODE)
	}

	ok, err := totp.ValidateCustom(*req.Code, *user.Secret, time.Now(), totp.ValidateOpts{
		Period:    myconstant.PASSCODE_EXPIRATION_DURATION,
		Digits:    6,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return response.InternalServerErrorErr(c, err)
	}

	if !ok {
		return response.Unauthorized(c, i18n.KEY_INTERNAL_ERROR_WRONG_PASSCODE)
	}

	EmailVerified := true
	dto := User{
		EmailVerified: &EmailVerified,
	}
	data, err := ctr.Repo.Update(user.ID, &dto)
	if err != nil {
		return response.InternalServerErrorErr(c, err)
	}

	userType := USER_TYPE_NORMAL
	newToken, err := ctr.RedisService.CreateToken(c.Request().Context(), data.ID, data.Email, &userType)
	if err != nil {
		return response.InternalServerErrorErr(c, err)
	}

	return response.Success(c, map[string]any{"accessToken": newToken, "user": data})
}

func (ctr *UserHandler) RequestEmailVerification(c echo.Context) error {
	type Req struct {
		Email *string `json:"email"`
	}
	var req Req
	if err := validate.Bind(c, &req); err != nil {
		return response.BadRequestErr(c, err)
	}

	user, err := ctr.Repo.GetOne(&UserRequestDto{ArrayEmail: &[]string{*req.Email}})
	if err != nil {
		if err == pgx.ErrNoRows {
			return response.BadRequest(c, i18n.KEY_NOT_FOUND_USER)
		}
		return response.BadRequestErr(c, err)
	}

	if user.EmailVerified != nil && *user.EmailVerified {
		return response.BadRequest(c, i18n.KEY_BAD_REQUEST_ALREADY_VERIFIED)
	}

	if user.Secret == nil {
		return response.InternalServerError(c, i18n.KEY_INTERNAL_ERROR_NO_SECRET_CODE)
	}

	checkEligibleToGetCode, releaseTime, err := ValidateGetPasscodeTime(user)
	if err != nil {
		return response.InternalServerErrorErr(c, err)
	}

	if !checkEligibleToGetCode {
		remainSec := int(time.Until(*releaseTime).Seconds())
		return c.JSON(http.StatusForbidden, map[string]any{
			"message":       i18n.Trsl(c, i18n.KEY_FORBIDDEN_WAIT_X_SECOND, strconv.Itoa(remainSec)),
			"statusCode":    http.StatusForbidden,
			"remain_second": remainSec,
			"release_time":  releaseTime,
		})
	}

	code, err := totp.GenerateCodeCustom(*user.Secret, time.Now(), totp.ValidateOpts{
		Period:    myconstant.PASSCODE_EXPIRATION_DURATION,
		Algorithm: otp.AlgorithmSHA1,
		Digits:    6,
	})
	if err != nil {
		return response.InternalServerErrorErr(c, err)
	}

	sendData := CreateEmailVerificationTemplate(c, user, &code)
	err = ctr.EmailService.Send(sendData)
	if err != nil {
		return response.InternalServerErrorErr(c, err)
	}

	currentTime := time.Now()
	dto := User{
		LastGetPasscodeTime: &currentTime,
	}
	_, err = ctr.Repo.Update(user.ID, &dto)
	if err != nil {
		return response.InternalServerErrorErr(c, err)
	}

	// to get duration in second
	sec := int(myconstant.PASSCODE_RESEND_DURATION.Seconds())
	return c.JSON(http.StatusOK, map[string]any{
		"statusCode":    http.StatusOK,
		"remain_second": sec,
		"release_time":  currentTime.Add(myconstant.PASSCODE_RESEND_DURATION),
	})
}

func (ctr *UserHandler) RequestPasswordResetEmail(c echo.Context) error {
	type Req struct {
		Email *string `json:"email"`
	}
	var req Req
	if err := validate.Bind(c, &req); err != nil {
		return response.BadRequestErr(c, err)
	}

	user, err := ctr.Repo.GetOne(&UserRequestDto{ArrayEmail: &[]string{*req.Email}})
	if err != nil {
		if err == pgx.ErrNoRows {
			return response.BadRequest(c, i18n.KEY_NOT_FOUND_USER)
		}
		return response.BadRequestErr(c, err)
	}

	if user.Secret == nil {
		return response.InternalServerError(c, i18n.KEY_INTERNAL_ERROR_NO_SECRET_CODE)
	}

	checkEligibleToGetCode, releaseTime, err := ValidateGetPasscodeTime(user)
	if err != nil {
		return response.InternalServerErrorErr(c, err)
	}

	if !checkEligibleToGetCode {
		remainSec := int(time.Until(*releaseTime).Seconds())
		return c.JSON(http.StatusForbidden, map[string]any{
			"message":       i18n.Trsl(c, i18n.KEY_FORBIDDEN_WAIT_X_SECOND, strconv.Itoa(remainSec)),
			"statusCode":    http.StatusForbidden,
			"remain_second": remainSec,
			"release_time":  releaseTime,
		})
	}

	code, err := totp.GenerateCodeCustom(*user.Secret, time.Now(), totp.ValidateOpts{
		Period:    myconstant.PASSCODE_EXPIRATION_DURATION,
		Algorithm: otp.AlgorithmSHA1,
		Digits:    6,
	})
	if err != nil {
		return response.InternalServerErrorErr(c, err)
	}

	sendData := CreateResetPasswordTemplate(c, user, &code)
	err = ctr.EmailService.Send(sendData)
	if err != nil {
		return response.InternalServerErrorErr(c, err)
	}

	currentTime := time.Now()
	dto := User{
		LastGetPasscodeTime: &currentTime,
	}
	_, err = ctr.Repo.Update(user.ID, &dto)
	if err != nil {
		return response.InternalServerErrorErr(c, err)
	}

	// to get duration in second
	sec := myconstant.PASSCODE_RESEND_DURATION / 1000000000
	return c.JSON(http.StatusOK, map[string]any{
		"statusCode":    http.StatusOK,
		"remain_second": sec,
		"release_time":  currentTime.Add(myconstant.PASSCODE_RESEND_DURATION),
	})
}
