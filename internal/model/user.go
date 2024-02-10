package model

import "errors"

type JwtTokenPair struct {
	AccessToken  *string `json:"access_token"`
	RefreshToken *string `json:"refresh_token"`
}

type User struct {
	ID             int    `json:"id" db:"id"`
	Email          string `json:"email" db:"email"`
	Username       string `json:"username" db:"username"`
	HashedPassword string `json:"password" db:"hashed_password"`
	CreatedAt
	UpdatedAt
}

type CreateUserInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignInUserInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserInput struct {
	Username *string `json:"username"`
}

func (o *UpdateUserInput) Validate() error {
	if o.Username == nil {
		return errors.New(validationError)
	}
	return nil
}
