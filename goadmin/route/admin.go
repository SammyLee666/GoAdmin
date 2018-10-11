package route

import (
	"github.com/gin-gonic/gin"
	"goadmin/goadmin/controllers"
)

func RegisterAdminARoute(r *gin.Engine){

	r.LoadHTMLGlob("goadmin/resources/views/**/*")

	admin := r.Group("/admin")
	{
		admin.GET("/", controllers.Index)
		admin.GET("/profile", controllers.UserProfile)

	}
}