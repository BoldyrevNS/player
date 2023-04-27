package DTO

type CategoryDTO struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type CreateCategoryDTO struct {
	Name string `json:"name"`
}
