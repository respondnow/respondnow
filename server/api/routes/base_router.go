package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/respondnow/respondnow/server/api/handlers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func BaseRouter(router *gin.RouterGroup) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/status", handlers.StatusHandler())
	router.GET("/version", handlers.InitVersionInfo())
}
