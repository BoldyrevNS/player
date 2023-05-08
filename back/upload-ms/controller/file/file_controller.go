package file

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shared/common/response"
	"upload-ms/DTO"
	"upload-ms/service/file"
)

type Controller interface {
	UploadVideo(ctx *gin.Context)
	GetVideo(ctx *gin.Context)
}

type controllerImpl struct {
	fileService file.Service
}

func NewUploadController(uploadService file.Service) Controller {
	return &controllerImpl{fileService: uploadService}
}

func (c *controllerImpl) UploadVideo(ctx *gin.Context) {
	var formData DTO.UploadFileDTO
	err := ctx.Bind(&formData)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = c.fileService.UploadVideo(formData)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
	}
}

func (c *controllerImpl) GetVideo(ctx *gin.Context) {
	var data DTO.GetFileDTO
	err := ctx.BindJSON(&data)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	fileObj, err := c.fileService.GetVideo(data)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	response.SendJSON(ctx, http.StatusOK, response.MessageJSON{
		Message: fileObj,
	})
}
