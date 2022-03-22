package routers

import (
    "dcard-project/controller"
    "dcard-project/controller/api/v1"
    _ "dcard-project/docs"
    "dcard-project/middleware"
    "fmt"
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

func SetRouter() *gin.Engine {
    r := gin.Default()
    r.Delims("{{{", "}}}")
    r.LoadHTMLGlob("resources/views/*")
    r.Static("/asset", "./resources/asset")
    r.GET("", controller.Index)

    url := ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", 8080))
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

    v1Group := r.Group("")
    v1Group.Use(middleware.IPLimitIntercept())
    {
        v1Group.POST("", v1.CreateUrl)
        v1Group.GET("/:shortUrl", v1.ToOrgPage)
    }

    return r
}
