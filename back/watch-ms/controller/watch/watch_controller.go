package watch

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shared/common/token"
	"watch-ms/DTO"
	watchService "watch-ms/service/watch"
)

type Controller interface {
	StartWatch(ctx *gin.Context)
}

type controllerImpl struct {
	watchService watchService.Service
}

func NewWatchController(watchService watchService.Service) Controller {
	return &controllerImpl{
		watchService: watchService,
	}
}

// StartWatch	godoc
// @Tags			Watch
// @Summary			User start watch title event
// @Security 		BearerAuth
// @Description 	Create new category in database
// @Param			category body DTO.UserStartWatchDTO true "Start watch info"
// @Success			201
// @Failure      	401
// @Failure      	500
// @Router			/watch/start [post]
func (c *controllerImpl) StartWatch(ctx *gin.Context) {
	var startWatchData DTO.UserStartWatchDTO
	err := ctx.ShouldBindJSON(&startWatchData)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	accessToken, err := token.GetTokenFromHeader(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	tokenData, err := token.ParseAccessToken(accessToken)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = c.watchService.UserStartWatch(startWatchData, tokenData.Id)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}
