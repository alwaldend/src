package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type CommonFields struct {
	ID        uuid.UUID    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at,omitempty"`
}

type CommonFieldsName struct {
	CommonFields
	Name string `gorm:"uniqueIndex;size:256" json:"name" validate:"required,min=3,max=256"`
}
