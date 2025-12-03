package response

import (
	"github.com/example/internal/i18n"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
)

func ErrHandling(c echo.Context, err error) string {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_username_key":
				return i18n.KEY_DUPLICATE_CONSTRAINT_USERNAME
			case "users_email_key":
				return i18n.KEY_DUPLICATE_CONSTRAINT_EMAIL
			case "community_community_identifier_key":
				return i18n.KEY_DUPLICATE_CONSTRAINT_COMMUNITY_IDENTIFIER
			default:
				return i18n.KEY_DUPLICATE_CONSTRAINT_UNIQUE_VALUE
			}
		}
	}
	return i18n.Trsl(c, err.Error())
}
