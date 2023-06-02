package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"shared/middlewares"
	uploadController "upload-ms/controller/upload"
	fileProvider "upload-ms/provider/file"
	"upload-ms/provider/watch_thumbnail"
	"upload-ms/service/episode"
	uploadService "upload-ms/service/file"
)

func episodesInit(router *gin.RouterGroup, dbInstance *gorm.DB) {
	saveFrameInfoCh := make(chan uploadService.SaveFrameInfo)
	uploadedFrameInfoCh := make(chan uploadService.UploadedFrameInfo)
	deleteFrameCh := make(chan string)

	protected := router.Group("/upload")
	protected.Use(middlewares.ProtectedMiddleware)

	adminPermission := router.Group("/upload")
	adminPermission.Use(middlewares.AdminPermissionMiddleware)

	fProvider := fileProvider.NewFileService()
	watchThumbnailProvider := watch_thumbnail.NewWatchThumbnailProvider(dbInstance)

	fileService := uploadService.NewFileService(fProvider, saveFrameInfoCh, uploadedFrameInfoCh, deleteFrameCh)
	episodeService := episode.NewEpisodeService(fileService, watchThumbnailProvider, saveFrameInfoCh, uploadedFrameInfoCh, deleteFrameCh)
	controller := uploadController.NewUploadController(episodeService)

	go episodeService.OnUserStartWatch()
	go fileService.SaveVideoFrameLocal()
	go fileService.DeleteVideoFrameLocal()
	go episodeService.SaveWatchThumbnail()

	protected.PATCH("/updateWatchThumbnail", controller.UpdateWatchThumbnail)

	adminPermission.POST("/episode", controller.UploadVideo)
}

func NewRouter(dbInstance *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	baseRouter := router.Group("/api/v1")
	episodesInit(baseRouter, dbInstance)
	return router
}
