package category

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"shared/common/response"
	"strconv"
	category2 "watch-ms/DTO"
	"watch-ms/service/category"
)

type Controller interface {
	GetAllCategories(ctx *gin.Context)
	CreateCategory(ctx *gin.Context)
	DeleteCategory(ctx *gin.Context)
}

type controllerImpl struct {
	service category.Service
}

func NewCategoryController(service category.Service) Controller {
	return &controllerImpl{
		service: service,
	}
}

// GetAllCategories godoc
// @Tags			Category
// @Summary			Get all categories
// @Security 		BearerAuth
// @Success			200 {object} response.DataJSON{data=[]DTO.CategoryDTO}
// @Failure      	401
// @Failure      	500
// @Router			/category/all [get]
func (c *controllerImpl) GetAllCategories(ctx *gin.Context) {
	categories, err := c.service.GetAllCategories()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	response.SendJSON(ctx, http.StatusOK, response.DataJSON{
		Data: categories,
	})
}

// CreateCategory	godoc
// @Tags			Category
// @Summary			Create new category
// @Security 		BearerAuth
// @Description 	Create new category in database
// @Param			category body DTO.CreateCategoryDTO true "Create category"
// @Success			201
// @Failure      	409 {object}  response.MessageJSON{}
// @Failure      	400
// @Failure      	401
// @Router			/category/ [post]
func (c *controllerImpl) CreateCategory(ctx *gin.Context) {
	var data category2.CreateCategoryDTO
	err := ctx.ShouldBindJSON(&data)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = c.service.CreateCategory(data)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.AbortWithStatus(http.StatusCreated)
}

// DeleteCategory 	godoc
// @Tags			Category
// @Summary			Delete category
// @Security 		BearerAuth
// @Description		Remove category data by id.
// @Param			categoryId   path   uint  true  "Category ID"
// @Success			200
// @Failure      	400 {object} response.MessageJSON{}
// @Failure      	500 {object} response.MessageJSON{}
// @Failure      	401
// @Router			/category/{categoryId} [delete]
func (c *controllerImpl) DeleteCategory(ctx *gin.Context) {
	categoryIdParam, find := ctx.Params.Get("categoryId")
	if !find {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.MessageJSON{
			Message: "provide category id",
		})
		return
	}
	categoryId, err := strconv.Atoi(categoryIdParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.MessageJSON{Message: "wrong format"})
		return
	}
	err = c.service.DeleteCategory(uint(categoryId))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.MessageJSON{Message: fmt.Sprintf("delete error: %v", err)})
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}
