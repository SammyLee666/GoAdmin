package route

import (
	"github.com/gin-gonic/gin"
	"goadmin/goadmin/controllers"
)

func RegisterAdminARoute(r *gin.Engine) {

	r.LoadHTMLGlob("goadmin/resources/views/**/*")

	admin := r.Group("/admin")
	{
		auth := admin.Group("/auth/menu")
		{
			auth.GET("", controllers.Show)
			auth.PUT("", controllers.Edit)
			auth.POST("", controllers.Post)
			auth.DELETE("", controllers.Dele)
		}
	}
}
