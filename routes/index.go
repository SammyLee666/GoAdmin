package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gin-contrib/sessions"
)

func Index(c *gin.Context) {

	//c.HTML(http.StatusOK, "home/index/general", gin.H{
	//	"Title":     "Sammy Blog",
	//	"Mooto":     "去过的地方越多，越知道自己想回到什么地方去。见过的人越多，越知道自己真正想待在什么人身边。",
	//	"MootoName": "夏正正",
	//})

	session := sessions.Default(c)
	msgInfo := session.Flashes("Info")
	session.Save()

	c.HTML(http.StatusOK, "goadmin/index/index", gin.H{
		"title":   "Login",
		"MsgInfo": msgInfo,
	})
}

func Test(c *gin.Context) {
	session := sessions.Default(c)
	session.AddFlash("Info flash.", "Info")
	session.Save()
}
