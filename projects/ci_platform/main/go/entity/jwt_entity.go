package entity

import "github.com/google/uuid"

type JWTClaims struct {
	ID    uuid.UUID `mapstructure:"id" validate:"required,min=1"`
	Admin bool      `mapstructure:"admin"`
}
