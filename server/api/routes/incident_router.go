package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/respondnow/respond/server/api/handlers"
)

func IncidentRouter(router *gin.RouterGroup) {
	router.POST("/create", handlers.CreateIncident())
	router.GET("/list", handlers.ListIncidents())
	router.GET("/get", handlers.GetIncident())
}
