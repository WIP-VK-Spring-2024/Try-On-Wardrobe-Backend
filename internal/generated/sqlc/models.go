// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package sqlc

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"
)

type Gender string

const (
	GenderMale    Gender = "male"
	GenderFemale  Gender = "female"
	GenderUnisex  Gender = "unisex"
	GenderUnknown Gender = "unknown"
)

func (e *Gender) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Gender(s)
	case string:
		*e = Gender(s)
	default:
		return fmt.Errorf("unsupported scan type for Gender: %T", src)
	}
	return nil
}

type NullGender struct {
	Gender Gender
	Valid  bool // Valid is true if Gender is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullGender) Scan(value interface{}) error {
	if value == nil {
		ns.Gender, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Gender.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullGender) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Gender), nil
}

type Season string

const (
	SeasonWinter Season = "winter"
	SeasonSpring Season = "spring"
	SeasonSummer Season = "summer"
	SeasonAutumn Season = "autumn"
)

func (e *Season) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Season(s)
	case string:
		*e = Season(s)
	default:
		return fmt.Errorf("unsupported scan type for Season: %T", src)
	}
	return nil
}

type NullSeason struct {
	Season Season
	Valid  bool // Valid is true if Season is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullSeason) Scan(value interface{}) error {
	if value == nil {
		ns.Season, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Season.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullSeason) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Season), nil
}

type Status string

const (
	StatusActive   Status = "active"
	StatusWishlist Status = "wishlist"
	StatusRepair   Status = "repair"
	StatusGiveAway Status = "give_away"
)

func (e *Status) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Status(s)
	case string:
		*e = Status(s)
	default:
		return fmt.Errorf("unsupported scan type for Status: %T", src)
	}
	return nil
}

type NullStatus struct {
	Status Status
	Valid  bool // Valid is true if Status is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullStatus) Scan(value interface{}) error {
	if value == nil {
		ns.Status, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Status.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Status), nil
}

type Clothes struct {
	ID        utils.UUID
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	Name      string
	Note      pgtype.Text
	UserID    utils.UUID
	StyleID   utils.UUID
	TypeID    utils.UUID
	SubtypeID utils.UUID
	Color     pgtype.Text
	Seasons   []domain.Season
	Image     string
	Status    NullStatus
}

type ClothesTag struct {
	ClothesID utils.UUID
	TagID     utils.UUID
}

type Outfit struct {
	ID         utils.UUID
	UserID     utils.UUID
	StyleID    utils.UUID
	CreatedAt  pgtype.Timestamptz
	UpdatedAt  pgtype.Timestamptz
	Name       pgtype.Text
	Note       pgtype.Text
	Image      pgtype.Text
	Transforms []byte
	Seasons    []domain.Season
}

type OutfitsTag struct {
	OutfitID utils.UUID
	TagID    utils.UUID
}

type Style struct {
	ID        utils.UUID
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	Name      string
}

type Subtype struct {
	ID        utils.UUID
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	Name      string
	TypeID    utils.UUID
}

type Tag struct {
	ID        utils.UUID
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	Name      string
	UseCount  int32
}

type TryOnResult struct {
	ID          utils.UUID
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	Rating      pgtype.Int4
	Image       string
	UserImageID utils.UUID
	ClothesID   utils.UUID
}

type Type struct {
	ID        utils.UUID
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	Name      string
	Tryonable bool
}

type User struct {
	ID        utils.UUID
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	Name      string
	Email     pgtype.Text
	Password  string
	Gender    NullGender
}

type UserImage struct {
	ID        utils.UUID
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	UserID    utils.UUID
	Image     string
}
