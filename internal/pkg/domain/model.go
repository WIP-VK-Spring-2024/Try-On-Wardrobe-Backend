package domain

import (
	"try-on/internal/pkg/utils"
)

type Model struct {
	ID utils.UUID `json:"uuid"`
	Timestamp
}

type Timestamp struct {
	CreatedAt utils.Time
	UpdatedAt utils.Time
}
