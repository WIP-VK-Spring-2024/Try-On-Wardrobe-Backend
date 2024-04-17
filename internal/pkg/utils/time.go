package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Time struct {
	time.Time
}

func (t Time) IsDefined() bool {
	return !t.IsZero()
}

func TimeFromPg(time pgtype.Timestamp) Time {
	return Time{Time: time.Time}
}

func TimeFromPgTz(time pgtype.Timestamptz) Time {
	return Time{Time: time.Time}
}
