package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/respondnow/respond/server/api/handlers"
	"github.com/respondnow/respond/server/api/middleware"
)

func AuthRouter(router *gin.RouterGroup) {
	router.POST("/signup", handlers.SignUp())
	router.POST("/login", handlers.Login())
	router.POST("/changePassword", handlers.ChangePassword())
	router.GET("/userMapping", middleware.AuthMiddleware(), handlers.GetUserMapping())
}
