package router

import (
	"auth-ms/controller"
	"auth-ms/provider"
	"auth-ms/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"shared/middlewares"
)

func authRoutes(router *gin.RouterGroup, dbInstance *gorm.DB) {
	public := router.Group("/auth")

	protected := router.Group("/auth")
	protected.Use(middlewares.ProtectedMiddleware)

	adminPermission := router.Group("/auth")
	adminPermission.Use(middlewares.AdminPermissionMiddleware)

	userProvider := provider.NewUserProvider(dbInstance)
	authService := service.NewAuthService(userProvider)
	authController := controller.NewAuthController(authService)

	public.POST("/", authController.Auth)
	public.POST("/registration", authController.Registration)
	public.POST("/refresh", authController.Refresh)

	protected.GET("/validateAuthToken", authController.ValidateAuthToken)

	adminPermission.DELETE(":userId", authController.DeleteUser)
	adminPermission.GET("/allUsers", authController.GetAllUsers)
}

func NewRouter(dbInstance *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	baseRouter := router.Group("/api/v1")
	authRoutes(baseRouter, dbInstance)

	return router
}
