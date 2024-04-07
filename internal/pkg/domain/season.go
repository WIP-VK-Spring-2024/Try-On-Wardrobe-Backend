package domain

type Season string

const (
	Winter Season = "winter"
	Spring Season = "spring"
	Summer Season = "summer"
	Autumn Season = "autumn"
)

var Seasons = []string{string(Winter), string(Spring), string(Summer), string(Autumn)}
