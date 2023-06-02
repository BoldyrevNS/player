package upload

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"upload-ms/DTO"
	"upload-ms/service/episode"
)

type Controller interface {
	UploadVideo(ctx *gin.Context)
	UpdateWatchThumbnail(ctx *gin.Context)
}

type controllerImpl struct {
	episodeService episode.Service
}

func NewUploadController(episodeService episode.Service) Controller {
	return &controllerImpl{
		episodeService: episodeService,
	}
}

// UploadVideo		godoc
// @Tags			Upload
// @Summary			Upload video
// @Param			uploadData formData DTO.UploadNewEpisodeDTO true "Upload video params"
// @Param			video formData file true "video file"
// @Param			thumbnail formData file true "episode thumbnail file"
// @Success			200
// @Failure      	400
// @Router			/upload/episode [post]
func (c *controllerImpl) UploadVideo(ctx *gin.Context) {
	var formData DTO.UploadNewEpisodeDTO
	err := ctx.Bind(&formData)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = c.episodeService.UploadNewEpisode(formData)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	ctx.AbortWithStatus(http.StatusOK)
}

// UpdateWatchThumbnail		godoc
// @Tags					Upload
// @Summary					Update thumbnail
// @Param					thumbnailData body DTO.UploadWatchThumbnailDTO true "Update existing watch thumbnail"
// @Success					200
// @Failure      			400
// @Failure      			500
// @Router					/upload/updateWatchThumbnail/ [patch]
func (c *controllerImpl) UpdateWatchThumbnail(ctx *gin.Context) {
	var data DTO.UploadWatchThumbnailDTO
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = c.episodeService.UploadWatchThumbnail(data)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}
