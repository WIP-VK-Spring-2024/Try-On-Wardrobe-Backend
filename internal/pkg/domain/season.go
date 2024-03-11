package domain

//go:generate stringer -type=Season
type Season uint8

const (
	Winter Season = iota
	Spring
	Summer
	Autumn
)
