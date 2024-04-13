package domain

import "try-on/internal/pkg/utils"

//easyjson:json
type Tag struct {
	Model
	Name     string
	UseCount int32
	EngName  string `json:"-"`
}

type TagUsecase interface {
	Get(limit, offset int) ([]Tag, error)
	Create(tag []string) error
	SetEngNames(tag []string) error
}

type TagRepository interface {
	Get(limit, offset int) ([]Tag, error)
	GetUserFavourite(userId utils.UUID, limit int) ([]Tag, error)
	GetNotCreated(tags []string) ([]string, error)
	Create(tags []Tag) error
	SetEngNames(tags, engNames []string) error
}
