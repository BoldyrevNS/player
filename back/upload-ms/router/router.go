package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
	"shared/broker"
	"shared/middlewares"
	uploadController "upload-ms/controller/file"
	fileProvider "upload-ms/provider"
	"upload-ms/service/episode"
	uploadService "upload-ms/service/file"
)

func episodesInit(router *gin.RouterGroup) {

	adminPermission := router.Group("/upload")
	adminPermission.Use(middlewares.AdminPermissionMiddleware)
	brokerWriter := broker.NewBrokerWriter("upload-video", os.Getenv("WRITE_BROKER"))
	fProvider := fileProvider.NewFileService()
	fileService := uploadService.NewFileService(fProvider)
	episodeService := episode.NewEpisodeService(fileService, brokerWriter)
	controller := uploadController.NewUploadController(episodeService)

	adminPermission.POST("/uploadEpisode", controller.UploadVideo)
}

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	baseRouter := router.Group("/api/v1")
	episodesInit(baseRouter)
	return router
}
