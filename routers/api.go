package routers

import (
	"dcard-project/app/Api/Controllers/UrlController"
	"dcard-project/app/Http/Controllers/IndexController"
	. "dcard-project/app/Http/Middleware"
	"dcard-project/app/Logic/Log"
	"github.com/gin-gonic/gin"
)

func SetRouter() {
	router := gin.Default()
	router.Delims("{{{", "}}}")
	router.LoadHTMLGlob("resources/views/*")
	router.Static("/asset", "./resources/asset")
	router.GET("", IndexController.Show)

	v1Router := router.Group("/link")
	v1Router.Use(IPLimitIntercept())
	{
		v1Router.POST("", UrlController.Create)
		v1Router.GET("/:shortUrl", UrlController.ToOrgPage)
	}

	err := router.Run()
	if err != nil {
		Log.Error.Println("Route Run Error")
		return
	}
}
