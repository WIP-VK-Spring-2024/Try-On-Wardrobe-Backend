package optional

import (
	"database/sql"

	"github.com/mailru/easyjson/jwriter"
)

type String struct {
	sql.NullString
}

func (s String) IsDefined() bool {
	return s.Valid
}

func (s String) MarshalEasyJSON(w *jwriter.Writer) {
	if s.Valid {
		w.String(s.String)
	} else {
		w.RawString("null")
	}
}

func (s *String) UnmarshalText(bytes []byte) error {
	if bytes != nil {
		s.Valid = true
		s.String = string(bytes)
	} else {
		s.Valid = false
	}
	return nil
}
