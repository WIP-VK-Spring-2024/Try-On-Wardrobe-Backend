package tags

import (
	"try-on/internal/pkg/domain"
	"try-on/internal/pkg/repository/sqlc/tags"
	"try-on/internal/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TagUsecase struct {
	translator domain.Translator
	repo       domain.TagRepository
}

func New(db *pgxpool.Pool, translator domain.Translator) domain.TagUsecase {
	return &TagUsecase{
		translator: translator,
		repo:       tags.New(db),
	}
}

func (t TagUsecase) Get(limit, offset int) ([]domain.Tag, error) {
	return t.repo.Get(limit, offset)
}

func (t TagUsecase) Create(tags []string) error {
	notCreated, err := t.repo.GetNotCreated(tags)
	if err != nil {
		return err
	}

	engNames, err := t.getEngNames(notCreated)
	if err != nil {
		return err
	}

	return t.repo.Create(utils.Zip(notCreated, engNames, func(name, engName string) domain.Tag {
		return domain.Tag{
			Name:    name,
			EngName: engName,
		}
	}))
}

func (t TagUsecase) SetEngNames(tags []string) error {
	engNames, err := t.getEngNames(tags)
	if err != nil {
		return err
	}
	return t.repo.SetEngNames(tags, engNames)
}

func (t TagUsecase) getEngNames(tags []string) ([]string, error) {
	engNames := make([]string, 0, len(tags))
	for _, tag := range tags {
		engName, err := t.translator.Translate(tag, domain.LanguageRU, domain.LanguageEN)
		if err != nil {
			return nil, err
		}
		engNames = append(engNames, engName)
	}
	return engNames, nil
}
