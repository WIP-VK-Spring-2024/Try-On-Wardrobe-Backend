package translate

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ToPgUUID(id uuid.UUID) pgtype.UUID {
	var res pgtype.UUID
	if id == uuid.Nil {
		return res
	}

	res.Bytes = id
	res.Valid = true

	return res
}
