package routers

import (
	"dcard-project/Api/Controllers/UrlControllers"
	"github.com/gin-gonic/gin"
)

func SetRouter() {
	router := gin.Default()

	v1Router := router.Group("/link")
	{
		v1Router.POST("", UrlControllers.Create)
		v1Router.GET("/:shortUrl", UrlControllers.ToOrgPage)
		//v1Router.GET("/:shortUrl", UrlControllers.Show)
	}

	router.Run()
}
