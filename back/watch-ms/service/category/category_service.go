package category

import (
	category2 "watch-ms/DTO"
	"watch-ms/model"
	"watch-ms/provider/category"
)

type Service interface {
	GetAllCategories() ([]category2.CategoryDTO, error)
	CreateCategory(data category2.CreateCategoryDTO) error
	DeleteCategory(data uint) error
}

type serviceImpl struct {
	provider category.Provider
}

func NewCategoryService(provider category.Provider) Service {
	return &serviceImpl{
		provider: provider,
	}
}

func (s *serviceImpl) GetAllCategories() ([]category2.CategoryDTO, error) {
	var categories []category2.CategoryDTO

	rawCategories, err := s.provider.GetAll()
	if err != nil {
		return categories, err
	}
	for _, cat := range rawCategories {
		categories = append(categories, category2.CategoryDTO{
			Id:   cat.Id,
			Name: cat.Name,
		})
	}
	return categories, nil
}

func (s *serviceImpl) CreateCategory(data category2.CreateCategoryDTO) error {
	cat := model.Category{
		Name: data.Name,
	}
	return s.provider.Create(cat)
}

func (s *serviceImpl) DeleteCategory(data uint) error {
	return s.provider.Delete(data)
}
