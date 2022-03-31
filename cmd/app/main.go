package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"dcard-project/controller"
	_ "dcard-project/docs"
	"dcard-project/pkg/postgres"
	"dcard-project/pkg/redis"
	"dcard-project/repository"
	"dcard-project/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	r := gin.Default()

	setupStatic(r)
	db := postgres.NewPgClient()
	client := redis.NewRedisClient()
	repo := repository.NewUrlRepo(db, client)
	urlSvc := service.NewUrlService(repo)
	controller.NewHandler(&controller.Config{R: r, UrlSvc: urlSvc})

	GracefulRunAndShutDown(r)
}

func GracefulRunAndShutDown(r *gin.Engine) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s \n", err)
		}
	}()

	<-ctx.Done()

	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server exiting")
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
