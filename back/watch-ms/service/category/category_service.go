package category

import (
	"watch-ms/DTO"
	"watch-ms/model"
	categoryProvider "watch-ms/provider/category"
)

type Service interface {
	GetAllCategories() ([]DTO.CategoryDTO, error)
	CreateCategory(data DTO.CreateCategoryDTO) error
	DeleteCategory(data uint) error
}

type serviceImpl struct {
	provider categoryProvider.Provider
}

func NewCategoryService(provider categoryProvider.Provider) Service {
	return &serviceImpl{
		provider: provider,
	}
}

func (s *serviceImpl) GetAllCategories() ([]DTO.CategoryDTO, error) {
	var categories []DTO.CategoryDTO

	rawCategories, err := s.provider.GetAll()
	if err != nil {
		return categories, err
	}
	for _, cat := range rawCategories {
		categories = append(categories, DTO.CategoryDTO{
			Id:   cat.Id,
			Name: cat.Name,
		})
	}
	return categories, nil
}

func (s *serviceImpl) CreateCategory(data DTO.CreateCategoryDTO) error {
	cat := model.Category{
		Name: data.Name,
	}
	return s.provider.Create(cat)
}

func (s *serviceImpl) DeleteCategory(data uint) error {
	return s.provider.Delete(data)
}
