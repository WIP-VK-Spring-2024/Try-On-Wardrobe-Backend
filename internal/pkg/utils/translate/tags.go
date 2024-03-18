package translate

import (
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/utils"
)

func TagsFromString(tags []string) []domain.Tag {
	return utils.Map(tags, func(tag *string) *domain.Tag {
		return &domain.Tag{Name: *tag}
	})
}

func TagsToString(tags []domain.Tag) []string {
	return utils.Map(tags, func(tag *domain.Tag) *string {
		return &tag.Name
	})
}
