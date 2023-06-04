package title

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shared/common/response"
	"shared/common/token"
	"watch-ms/DTO"
	titleService "watch-ms/service/title"
)

type Controller interface {
	CreateTitle(ctx *gin.Context)
	GetUserTitles(ctx *gin.Context)
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

// GetUserTitles godoc
// @Tags			Title
// @Summary			Get all user titles exclude watched
// @Security 		BearerAuth
// @Success			200 {object} response.DataJSON{data=[]DTO.TitleDTO}
// @Failure      	401
// @Failure      	500
// @Router			/title/user [get]
func (c *controllerImpl) GetUserTitles(ctx *gin.Context) {
	headerToken, err := token.GetTokenFromHeader(ctx)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	tokenData, err := token.ParseAccessToken(headerToken)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	titles, err := c.titleService.GetTitlesExcludeWatched(tokenData.Id)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	response.SendJSON(ctx, http.StatusOK, response.DataJSON{
		Data: titles,
	})
}
