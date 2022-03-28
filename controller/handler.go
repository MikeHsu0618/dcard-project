package controller

import (
    "dcard-project/middleware"
    "dcard-project/model"
    "github.com/gin-gonic/gin"
)

type Handler struct {
    svc model.UrlService
}

func NewHandler(r *gin.Engine, s model.UrlService) {
    handler := &Handler{
        svc: s,
    }

    v1Group := r.Group("/")
    v1Group.Use(middleware.IPLimitIntercept())
    {
        v1Group.POST("", handler.svc.CreateUrl)
        v1Group.GET(":shortUrl", handler.svc.ToOrgPage)
    }
}
