package main

import (
	"github.com/gin-gonic/gin"
	"os"
	_ "github.com/joho/godotenv/autoload"
	"goadmin/routes"
	"goadmin/goadmin/route"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/mongo"
	"goadmin/db"
	"html/template"
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

	r.SetFuncMap(template.FuncMap{
		"OutputHTML": func(html string) template.HTML {
			return template.HTML(html)
		},
		"OutputJS": func(js string) template.JS {
			return template.JS(js)
		},
	})

	//templates
	r.LoadHTMLGlob("templates/**/**/*")

	// session
	store := mongo.NewStore(db.Sessions, 3600, true, []byte("secret"))
	//store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{
		//Domain: "localhost",
		MaxAge: 3 * 24 * 3600,
	})
	r.Use(sessions.Sessions("goulang", store))
}

func loadRouters(r *gin.Engine) {
	r.GET("", routes.Index)
	r.GET("test", routes.Test)

}
