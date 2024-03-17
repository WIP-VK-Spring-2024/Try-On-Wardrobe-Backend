package usecase

import (
	"database/sql"

	"try-on/internal/pkg/domain"

	"github.com/google/uuid"
)

type ClothesUsecase struct {
	repo domain.ClothesRepository
}

func New(repo domain.ClothesRepository) domain.ClothesUsecase {
	return &ClothesUsecase{
		repo: repo,
	}
}

func (c *ClothesUsecase) Create(clothes *domain.Clothes) error {
	return c.repo.Create(toModel(clothes))
}

func (c *ClothesUsecase) Update(clothes *domain.Clothes) error {
	return c.repo.Update(toModel(clothes))
}

func (c *ClothesUsecase) Get(id uuid.UUID) (*domain.Clothes, error) {
	clothesModel, err := c.repo.Get(id)
	if err != nil {
		return nil, err
	}

	return fromModel(clothesModel), nil
}

func (c *ClothesUsecase) Delete(id uuid.UUID) error {
	return c.repo.Delete(id)
}

func (c *ClothesUsecase) GetByUser(userID uuid.UUID, limit int) ([]domain.Clothes, error) {
	clothes, err := c.repo.GetByUser(userID, limit)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Clothes, 0, len(clothes))
	for _, clothing := range clothes {
		result = append(result, *fromModel(&clothing))
	}

	return result, nil
}

func fromModel(clothesModel *domain.ClothesModel) *domain.Clothes {
	clothes := &domain.Clothes{
		ID:      clothesModel.ID,
		UserID:  clothesModel.UserID,
		Image:   clothesModel.Image,
		Name:    clothesModel.Name,
		Type:    clothesModel.Type.Name,
		Subtype: clothesModel.Subtype.Name,
		Color:   clothesModel.Color.String,
		Seasons: clothesModel.Seasons,
	}

	if clothesModel.Style != nil {
		clothes.Style = clothesModel.Style.Name
	}

	clothes.Tags = make([]string, 0, len(clothesModel.Tags))
	for _, tag := range clothesModel.Tags {
		clothes.Tags = append(clothes.Tags, tag.Name)
	}

	return clothes
}

func toModel(clothes *domain.Clothes) *domain.ClothesModel {
	model := &domain.ClothesModel{
		UserID: clothes.UserID,
		Name:   clothes.Name,
		Type: domain.Type{
			Name: clothes.Type,
		},
		Subtype: domain.Subtype{
			Name: clothes.Subtype,
		},
		Image: clothes.Image,
	}

	model.Type.ID = uuid.Nil
	model.Subtype.ID = uuid.Nil

	if clothes.Color != "" {
		model.Color = sql.NullString{String: clothes.Color, Valid: true}
	}

	if clothes.Note != "" {
		model.Note = sql.NullString{String: clothes.Note, Valid: true}
	}

	if clothes.Style != "" {
		model.Style = &domain.Style{Name: clothes.Style}
		model.Style.ID = uuid.Nil
	}

	model.Tags = make([]domain.Tag, 0, len(clothes.Tags))
	for _, tag := range clothes.Tags {
		model.Tags = append(model.Tags, domain.Tag{Name: tag})
	}

	return model
}
