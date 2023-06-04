package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"shared/middlewares"
	"watch-ms/DI"
	"watch-ms/controller/category"
	"watch-ms/controller/season"
	"watch-ms/controller/title"
	"watch-ms/controller/watch"
)

func categoryInit(router *gin.RouterGroup, controller category.Controller) {
	protected := router.Group("/category")
	protected.Use(middlewares.ProtectedMiddleware)

	adminPermission := router.Group("/category")
	adminPermission.Use(middlewares.AdminPermissionMiddleware)

	protected.GET("/all", controller.GetAllCategories)

	adminPermission.POST("/", controller.CreateCategory)
	adminPermission.DELETE("/:categoryId", controller.DeleteCategory)
}

func watchInit(router *gin.RouterGroup, controller watch.Controller) {
	protected := router.Group("/watch")
	protected.Use(middlewares.ProtectedMiddleware)

	protected.POST("/start", controller.StartWatch)
}

func titleInit(router *gin.RouterGroup, controller title.Controller) {
	adminPermission := router.Group("/title")
	adminPermission.Use(middlewares.AdminPermissionMiddleware)

	protected := router.Group("/title")
	protected.Use(middlewares.ProtectedMiddleware)

	adminPermission.POST("/", controller.CreateTitle)

	protected.GET("/user", controller.GetUserTitles)
}

func seasonInit(router *gin.RouterGroup, controller season.Controller) {
	adminPermission := router.Group("/season")
	adminPermission.Use(middlewares.AdminPermissionMiddleware)

	protected := router.Group("/season")
	protected.Use(middlewares.ProtectedMiddleware)

	adminPermission.POST("/", controller.CreateSeason)

	protected.GET("/title", controller.GetAllTitleSeasons)
}

func NewRouter(container DI.ControllerContainer) *gin.Engine {
	router := gin.Default()
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	baseRouter := router.Group("/api/v1")
	categoryInit(baseRouter, container.CategoryController)
	watchInit(baseRouter, container.WatchController)
	titleInit(baseRouter, container.TitleController)
	seasonInit(baseRouter, container.SeasonController)

	return router
}
