package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {

	c.HTML(http.StatusOK, "home/index/general", gin.H{
		"Title":     "Sammy Blog",
		"Mooto":     "去过的地方越多，越知道自己想回到什么地方去。见过的人越多，越知道自己真正想待在什么人身边。",
		"MootoName": "夏正正",
	})
}
