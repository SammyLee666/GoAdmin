package main

import (
	"github.com/gin-gonic/gin"
	"os"
	_ "github.com/joho/godotenv/autoload"
	"goadmin/routes"
	"goadmin/goadmin/route"
	"github.com/gin-contrib/sessions/mongo"
	"github.com/gin-contrib/sessions"
	"goadmin/db"
)

func main() {
	router := gin.New()

	loadMiddlewares(router)
	loadRouters(router)
	//goadmin
	route.RegisterAdminARoute(router)
	router.Run(":" + os.Getenv("PORT"))
}

func loadMiddlewares(r *gin.Engine) {
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	//admin 静态资源
	r.Static("/assets", "./goadmin/resources/assets")

	r.LoadHTMLGlob("templates/**/**/*")

	// session
	store := mongo.NewStore(db.Sessions, 3600, true, []byte("secret"))
	r.Use(sessions.Sessions("goulang", store))
}

func loadRouters(r *gin.Engine) {
	r.GET("", routes.Index)


}
