package gender

//go:generate stringer -type=Gender
type Gender uint8

const (
	Male Gender = iota
	Female
	Unisex
)
