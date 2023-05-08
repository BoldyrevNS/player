package DTO

import "mime/multipart"

type UploadFileDTO struct {
	File      *multipart.FileHeader `form:"file"`
	TitleName string                `form:"title_name"`
	TitleId   uint                  `form:"title_id"`
}

type GetFileDTO struct {
	Filename string `json:"filename"`
	Title    string `json:"title"`
}
