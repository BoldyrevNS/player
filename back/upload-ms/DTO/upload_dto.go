package DTO

import "mime/multipart"

type UploadNewEpisodeDTO struct {
	Video         *multipart.FileHeader `form:"video"  swaggerignore:"true"`
	Thumbnail     *multipart.FileHeader `form:"thumbnail" swaggerignore:"true"`
	SeasonId      uint                  `form:"seasonId"`
	EpisodeName   string                `form:"episodeName"`
	EpisodeNumber uint                  `form:"episodeNumber"`
}

type UploadWatchThumbnailDTO struct {
	TimeCode string `form:"timeCode"`
	WatchId  uint   `form:"watchId"`
}
