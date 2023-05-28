package events

type UploadNewEpisodeDTO struct {
	VideoPath     string `json:"videoPath"`
	ThumbnailPath string `json:"thumbnailPath"`
	SeasonId      uint   `json:"seasonId"`
	EpisodeName   string `json:"episodeName"`
	EpisodeNumber uint   `json:"episodeNumber"`
}
