package utils

import (
	"time"
)

type Time struct {
	time.Time
}

func (t Time) IsDefined() bool {
	return !t.IsZero()
}
