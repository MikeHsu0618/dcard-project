package routers

import (
	"dcard-project/app/Api/Controllers/UrlControllers"
	"dcard-project/app/Http/Controllers/IndexController"
	"github.com/gin-gonic/gin"
)

func SetRouter() {
	router := gin.Default()
	//config := cors.DefaultConfig()
	//config.AllowAllOrigins = true
	//router.Use(cors.New(config))
	router.Delims("{{{", "}}}")
	router.LoadHTMLGlob("resources/views/*")
	router.Static("/asset", "./resources/asset")

	router.GET("", IndexController.Show)

	v1Router := router.Group("/link")
	{
		v1Router.POST("", UrlControllers.Create)
		v1Router.GET("/:shortUrl", UrlControllers.ToOrgPage)
	}

	router.Run()
}
