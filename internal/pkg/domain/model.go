package domain

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	AutoTimestamp
}

type AutoTimestamp struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}
