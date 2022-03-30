package main

import (
	"dcard-project/controller"
	_ "dcard-project/docs"
	"dcard-project/pkg/postgres"
	"dcard-project/pkg/redis"
	"dcard-project/repository"
	"dcard-project/service"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func main() {
	r := gin.Default()

	setupStatic(r)
	db := postgres.NewPgClient()
	client := redis.NewRedisClient()
	repo := repository.NewUrlRepo(db, client)
	urlSvc := service.NewUrlService(repo)
	controller.NewHandler(&controller.Config{R: r, UrlSvc: urlSvc})

	r.Run()
}

func setupStatic(r *gin.Engine) {
	r.Delims("{{{", "}}}")
	r.LoadHTMLGlob("resources/views/*")
	r.Static("/asset", "./resources/asset")

	url := ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", 8080))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	r.GET("", func(c *gin.Context) {
		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{"title": "短網址產生器"},
		)
	})
}
