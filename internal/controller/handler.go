package controller

import (
	"dcard-project/internal/middleware"
	"dcard-project/internal/repository"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	urlSvc repository.UrlService
}

type Config struct {
	R      *gin.Engine
	UrlSvc repository.UrlService
}

func NewHandler(c *Config) {
	handler := &Handler{
		urlSvc: c.UrlSvc,
	}

	v1Group := c.R.Group("/")
	v1Group.Use(middleware.IPLimitIntercept())
	{
		v1Group.POST("", handler.urlSvc.CreateUrl)
		v1Group.GET(":shortUrl", handler.urlSvc.ToOrgPage)
	}
}
