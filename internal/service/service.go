package service

import (
	"my_app/internal/model"
	"my_app/internal/repository"
)

type Authorization interface {
	CreateUser(input model.CreateUserInput) (int, error)
	UserAuthorize(input model.SignInUserInput) (int, error)
	GenerateJwtTokenPair(userId int) (model.JwtTokenPair, error)
	GetUserIdFromAuthHeader(token string) (int, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{Authorization: NewAuthorizationService(repos)}
}
