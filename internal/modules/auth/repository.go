package auth

import (
	"github.com/example/internal/models"
	"github.com/example/internal/pkg/sqlhelper"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	DB *pgxpool.Pool
}

func (repo *AuthRepository) checkEmailExist(email *string) (*bool, error) {
	parameters := []any{email}
	query := `SELECT EXISTS (
				SELECT 1
				FROM users
				WHERE email = $1
	`

	filter := []string{
		models.COMMON_STATUS_OPEN,
		models.COMMON_STATUS_CLOSED,
		models.COMMON_STATUS_PENDING,
	}
	query += sqlhelper.ArrayContain("status", &parameters, &filter, "AND", nil)

	query += " ); "
	var check bool
	err := sqlhelper.Check(repo.DB, &check, query, parameters...)
	if err != nil {
		return nil, err
	}

	return &check, nil
}
