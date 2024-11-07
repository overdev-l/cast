package api

import (
	api "cast/internal/api/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/point", api.PointController)
	return router
}
