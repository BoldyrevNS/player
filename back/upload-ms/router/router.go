package router

import (
	"github.com/gin-gonic/gin"
	upload_controller "upload-ms/controller/file"
	fileProvider "upload-ms/provider"
	upload_service "upload-ms/service/file"
)

func uploadFileRoutes(router *gin.RouterGroup) {

	adminPermission := router.Group("/upload")
	fProvider := fileProvider.NewFileService()
	service := upload_service.NewUploadService(fProvider)
	controller := upload_controller.NewUploadController(service)

	adminPermission.POST("/", controller.UploadVideo)
	adminPermission.POST("/video", controller.GetVideo)
}

func NewRouter() *gin.Engine {
	router := gin.Default()

	baseRouter := router.Group("/api/v1")
	uploadFileRoutes(baseRouter)
	return router
}
