package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func MetricRouter(router *gin.RouterGroup) {
	router.GET("/", gin.WrapH(promhttp.Handler()))
}
