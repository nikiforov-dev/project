package repository

import (
	"github.com/jmoiron/sqlx"
	"my_app/internal/model"
)

type User interface {
	CreateUser(input model.CreateUserInput) (int, error)
	GetUserById(id int) (model.User, error)
	GetUserByLogin(login string) (model.User, error)
	UpdateUserById(id int, input model.UpdateUserInput) error
}

type Repository struct {
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUserRepoPostgres(db),
	}
}
