package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/respondnow/respond/server/api/handlers"
)

func AuthRouter(router *gin.RouterGroup) {
	router.POST("/signup", handlers.SignUp())
	router.POST("/login", handlers.Login())
	router.POST("/changePassword", handlers.ChangePassword())
}
