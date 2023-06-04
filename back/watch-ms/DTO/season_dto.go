package DTO

type CreateSeasonDTO struct {
	TitleId uint `json:"titleId"`
	Number  uint `json:"number"`
}

type GetTitleSeasonsDTO struct {
	TitleId uint `json:"titleId"`
}

type TitleSeasonDTO struct {
	SeasonId uint `json:"seasonId"`
	Number   uint `json:"number"`
}
