package file

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"upload-ms/DTO"
	"upload-ms/service/episode"
)

type Controller interface {
	UploadVideo(ctx *gin.Context)
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
// @Router			/upload/ [post]
func (c *controllerImpl) UploadVideo(ctx *gin.Context) {
	var formData DTO.UploadNewEpisodeDTO
	err := ctx.Bind(&formData)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = c.episodeService.UploadNewEpisode(formData)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
	ctx.AbortWithStatus(http.StatusOK)
}
