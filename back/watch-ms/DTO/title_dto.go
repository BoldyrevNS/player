package DTO

type CreateTitleDTO struct {
	CategoryId uint   `json:"categoryId"`
	Name       string `json:"name"`
}

type TitleDTO struct {
	Id           uint   `json:"id"`
	Name         string `json:"name"`
	ThumbnailUrl string `json:"thumbnailUrl"`
	CategoryName string `json:"categoryName"`
	CategoryId   uint   `json:"categoryId"`
}
