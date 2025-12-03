package response

import (
	"net/http"

	"github.com/example/internal/i18n"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func Error(c echo.Context, code int, err error) error {
	return c.JSON(code, ErrorResponse{Message: ErrHandling(c, err), StatusCode: code})
}

func Forbidden(c echo.Context, args ...string) error {
	return c.JSON(http.StatusForbidden, ErrorResponse{Message: i18n.Trsl(c, args...), StatusCode: http.StatusForbidden})
}

func NotFound(c echo.Context, args ...string) error {
	return c.JSON(http.StatusNotFound, ErrorResponse{Message: i18n.Trsl(c, args...), StatusCode: http.StatusNotFound})
}

func Unauthorized(c echo.Context, args ...string) error {
	return c.JSON(http.StatusUnauthorized, ErrorResponse{Message: i18n.Trsl(c, args...), StatusCode: http.StatusBadRequest})
}

func BadRequest(c echo.Context, args ...string) error {
	return c.JSON(http.StatusBadRequest, ErrorResponse{Message: i18n.Trsl(c, args...), StatusCode: http.StatusBadRequest})
}

func BadRequestErr(c echo.Context, err error) error {
	code := http.StatusBadRequest
	if err == pgx.ErrNoRows {
		code = http.StatusNotFound
	}
	return c.JSON(code, ErrorResponse{Message: ErrHandling(c, err), StatusCode: code})
}

func InternalServerError(c echo.Context, args ...string) error {
	return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: i18n.Trsl(c, args...), StatusCode: http.StatusInternalServerError})
}

func InternalServerErrorErr(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: ErrHandling(c, err), StatusCode: http.StatusInternalServerError})
}

func Success(c echo.Context, message ...any) error {
	if len(message) == 0 {
		return c.JSON(http.StatusOK, "OK")
	}
	return c.JSON(http.StatusOK, message[0])
}
