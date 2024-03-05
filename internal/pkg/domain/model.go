package domain

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID uuid.UUID
	AutoTimestamp
}

type AutoTimestamp struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}
