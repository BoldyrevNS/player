package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"shared/middlewares"
	categoryController "watch-ms/controller/category"
	categoryProvider "watch-ms/provider/category"
	episodeProvider "watch-ms/provider/episode"
	categoryService "watch-ms/service/category"
	"watch-ms/service/episode"
)

func categoryInit(router *gin.RouterGroup, dbInstance *gorm.DB) {
	protected := router.Group("/category")
	protected.Use(middlewares.ProtectedMiddleware)

	adminPermission := router.Group("/category")
	adminPermission.Use(middlewares.AdminPermissionMiddleware)

	provider := categoryProvider.NewCategoryProvider(dbInstance)
	service := categoryService.NewCategoryService(provider)
	controller := categoryController.NewCategoryController(service)

	protected.GET("/all", controller.GetAllCategories)

	adminPermission.POST("/", controller.CreateCategory)
	adminPermission.DELETE("/:categoryId", controller.DeleteCategory)
}

func episodeInit(dbInstance *gorm.DB) {
	provider := episodeProvider.NewEpisodeProvider(dbInstance)
	episodeService := episode.NewEpisodeService(provider)
	go episodeService.UploadNewEpisode()
}

func NewRouter(dbInstance *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	baseRouter := router.Group("/api/v1")
	categoryInit(baseRouter, dbInstance)
	episodeInit(dbInstance)

	return router
}
