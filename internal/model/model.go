package model

import "time"

const (
	validationError = "validation error"
)

type CreatedAt struct {
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
}

type UpdatedAt struct {
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}
