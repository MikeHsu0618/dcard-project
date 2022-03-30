package controller

import (
	"dcard-project/middleware"
	"dcard-project/model"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	urlSvc model.UrlService
}

type Config struct {
	R      *gin.Engine
	UrlSvc model.UrlService
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
