package user

import (
	"errors"

	"github.com/example/internal/i18n"
	"github.com/example/internal/models"
	"github.com/example/internal/pkg/sqlhelper"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func (repo *UserRepository) Create(dto *User) (*User, error) {
	data, err := sqlhelper.Create(repo.DB, dto)
	if err != nil {
		return nil, err
	}

	if len(*data) == 0 {
		return nil, errors.New(i18n.KEY_INTERNAL_ERROR_USER_NOT_CREATED)
	}
	return &(*data)[0], nil
}

func (repo *UserRepository) GetAll(dto *UserRequestDto) (*[]User, error) {
	parameters := []any{}

	query := repo.selectQuery(dto, &parameters)
	query += sqlhelper.DefaultOrder("")
	query += sqlhelper.Pagination(dto.Page, dto.Limit)

	users, err := sqlhelper.Query[User](repo.DB, query, parameters...)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (repo *UserRepository) GetOne(dto *UserRequestDto) (*User, error) {
	parameters := []any{}

	query := repo.selectQuery(dto, &parameters)

	query += "LIMIT 1"

	users, err := sqlhelper.Query[User](repo.DB, query, parameters...)
	if err != nil {
		return nil, err
	}

	if len(*users) == 0 {
		return nil, pgx.ErrNoRows
	}

	return &(*users)[0], nil
}

func (repo *UserRepository) Update(userID *string, data *User) (*User, error) {
	err := sqlhelper.Update(repo.DB, data, "id = ?", userID)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (repo *UserRepository) selectQuery(dto *UserRequestDto, parameters *[]any) string {
	query := `SELECT * from users WHERE id IS NOT NULL  `

	query += sqlhelper.ArrayContain("email", parameters, dto.ArrayEmail, "AND", nil)
	query += sqlhelper.ArrayContain("id", parameters, dto.ArrayID, "AND", nil)
	query += sqlhelper.ArrayContain("status", parameters, dto.ArrayStatus, "AND", []string{models.COMMON_STATUS_OPEN})

	return query
}
