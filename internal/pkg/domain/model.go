package domain

import (
	"try-on/internal/pkg/utils"
)

type Model struct {
	ID utils.UUID `json:"uuid"`
	AutoTimestamp
}

type AutoTimestamp struct {
	CreatedAt utils.Time
	UpdatedAt utils.Time `json:"-"`
}
