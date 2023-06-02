package events

type UserStartWatchDTO struct {
	WatchId  uint   `json:"userId"`
	VideoUrl string `json:"videoUrl"`
}

type UpdateThumbnailWatchDTO struct {
	WatchId       uint   `json:"watchId"`
	ThumbnailPath string `json:"thumbnailPath"`
}
