package DTO

import "mime/multipart"

type UploadNewEpisodeDTO struct {
	Video         *multipart.FileHeader `form:"video"  swaggerignore:"true"`
	Thumbnail     *multipart.FileHeader `form:"thumbnail" swaggerignore:"true"`
	SeasonId      uint                  `form:"seasonId"`
	TitleId       uint                  `form:"titleId"`
	EpisodeName   string                `form:"episodeName"`
	EpisodeNumber uint                  `form:"episodeNumber"`
}

type UploadWatchThumbnailDTO struct {
	Thumbnail *multipart.FileHeader `form:"thumbnail" swaggerignore:"true"`
	WatchId   uint                  `form:"watch_id"`
}
