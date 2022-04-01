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
	"dcard-project/pkg/logger"
	"dcard-project/pkg/postgres"
	"dcard-project/pkg/redis"
	"dcard-project/repository"
	"dcard-project/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	r  *gin.Engine
	lg *logger.Logger
)

func main() {
	r = gin.Default()
	setupStatic(r)
	db := postgres.NewPgClient()
	client := redis.NewRedisClient()
	lg = logger.NewLogger()

	repo := repository.NewUrlRepo(db, client, lg)
	urlSvc := service.NewUrlService(repo)
	controller.NewHandler(&controller.Config{R: r, UrlSvc: urlSvc})

	GracefulRunAndShutDown()
}

func GracefulRunAndShutDown() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			lg.Fatal("listen: " + err.Error())
		}
	}()

	<-ctx.Done()

	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		lg.Fatal("Server forced to shutdown: " + err.Error())
	}
	lg.Info("Server exiting")
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
