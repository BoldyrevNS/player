package season

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shared/common/response"
	"watch-ms/DTO"
	"watch-ms/service/season"
)

type Controller interface {
	CreateSeason(ctx *gin.Context)
	GetAllTitleSeasons(ctx *gin.Context)
}

type controllerImpl struct {
	seasonService season.Service
}

func NewSeasonController(seasonService season.Service) Controller {
	return &controllerImpl{
		seasonService: seasonService,
	}
}

// CreateSeason	godoc
// @Tags			Season
// @Summary			Create new season
// @Security 		BearerAuth
// @Description 	Create new season in database
// @Param			category body DTO.CreateSeasonDTO true "Create season"
// @Success			201
// @Failure      	400
// @Failure      	401
// @Router			/season/ [post]
func (c *controllerImpl) CreateSeason(ctx *gin.Context) {
	var data DTO.CreateSeasonDTO
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = c.seasonService.CreateSeason(data)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.AbortWithStatus(http.StatusCreated)
}

// GetAllTitleSeasons 	godoc
// @Tags			  	Season
// @Summary				Get all seasons of title
// @Security 			BearerAuth
// @Success				200 {object} response.DataJSON{data=[]DTO.TitleSeasonDTO}
// @Failure      		401
// @Failure      		500
// @Router				/season/title [get]
func (c *controllerImpl) GetAllTitleSeasons(ctx *gin.Context) {
	var data DTO.GetTitleSeasonsDTO
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	seasons, err := c.seasonService.GetAllTitleSeasons(data)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	response.SendJSON(ctx, http.StatusOK, response.DataJSON{
		Data: seasons,
	})
}
