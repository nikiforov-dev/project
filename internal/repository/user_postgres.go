package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"my_app/internal/model"
	"strings"
)

type UserRepoPostgres struct {
	db *sqlx.DB
}

func NewUserRepoPostgres(db *sqlx.DB) *UserRepoPostgres {
	return &UserRepoPostgres{db: db}
}

func (r *UserRepoPostgres) CreateUser(user model.CreateUserInput) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (email, username, hashed_password) VALUES ($1, $2, $3) RETURNING id", usersTable)

	res := r.db.QueryRow(
		query,
		user.Email,
		user.Username,
		user.Password,
	)
	if err := res.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserRepoPostgres) GetUserById(id int) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT id, email, username, hashed_password, access_token, refresh_token FROM %s WHERE id = $1", usersTable)

	if err := r.db.Get(&user, query, id); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *UserRepoPostgres) GetUserByLogin(login string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT id, username, email, hashed_password, created_at, updated_at FROM %s WHERE email = $1 OR username = $1", usersTable)

	if err := r.db.Get(&user, query, login); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *UserRepoPostgres) UpdateUserById(id int, input model.UpdateUserInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Username != nil {
		setValues = append(setValues, fmt.Sprintf("username=$%d", argId))
		args = append(args, *input.Username)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	args = append(args, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id=$%s", usersTable, setQuery, argId+1)

	_, err := r.db.Exec(query, args...)
	return err
}
