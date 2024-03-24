package domain

import (
	"time"

	"try-on/internal/pkg/utils"
)

type Model struct {
	ID utils.UUID `json:"uuid"`
	AutoTimestamp
}

type AutoTimestamp struct {
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
