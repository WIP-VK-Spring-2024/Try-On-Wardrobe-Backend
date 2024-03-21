package utils

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

type UUID struct {
	uuid.UUID
}

func (u UUID) IsDefined() bool {
	return u.UUID != uuid.Nil
}

func ParseUUID(str string) (UUID, error) {
	u, err := uuid.Parse(str)
	return UUID{UUID: u}, err
}

func (u *UUID) Scan(value any) error {
	return u.UUID.Scan(value)
}

func (u UUID) Value() (driver.Value, error) {
	if !u.IsDefined() {
		return nil, nil
	}
	return u.UUID.Value()
}

var NilUUID UUID = UUID{
	UUID: uuid.Nil,
}
