package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Point struct {
	LiveId int    `form:"liveId"`
	Plat   string `form:"plat"`
	Url    string `form:"url"`
}

func PointController(c *gin.Context) {
	var pointBody Point
	if err := c.ShouldBind(&pointBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error",
		})
	}
}
