package title

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"watch-ms/DTO"
	titleService "watch-ms/service/title"
)

type Controller interface {
	CreateTitle(ctx *gin.Context)
}

type controllerImpl struct {
	titleService titleService.Service
}

func NewTitleController(titleService titleService.Service) Controller {
	return &controllerImpl{
		titleService: titleService,
	}
}

// CreateTitle	godoc
// @Tags			Title
// @Summary			Create new title
// @Security 		BearerAuth
// @Description 	Create new title in database
// @Param			category body DTO.CreateTitleDTO true "Create title"
// @Success			201
// @Failure      	400
// @Failure      	401
// @Router			/title/ [post]
func (c *controllerImpl) CreateTitle(ctx *gin.Context) {
	var data DTO.CreateTitleDTO
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = c.titleService.CreateNewTitle(data)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.AbortWithStatus(http.StatusCreated)
}
