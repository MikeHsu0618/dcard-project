package model

import (
    "dcard-project/pkg/goquery"
    "github.com/gin-gonic/gin"
)

type Url struct {
    ID     int64  `json:"id" form:"id"`                              // 列名 `id`
    OrgUrl string `json:"org_url" form:"org_url" binding:"required"` // 列名 `org_url`
}

type ApiUrl struct {
    ShortUrl string `json:"short_url" form:"short_url" example:"5Cb"`
    goquery.Meta
}

type CreateUrl struct {
    OrgUrl string `json:"org_url" example:"https://www.google.com"`
}

type UrlRepository interface {
    Create(url *Url) (err error)
    GetById(urlId int64, url *Url) (err error)
    GetByOrgUrl(orgUrl string) (url *Url, err error)
    Lock(key string) bool
    UnLock(key string) int64
    GetCache(key int64) (result string, err error)
    SetCache(key int64, orgUrl string)
}

type UrlService interface {
    CreateUrl(c *gin.Context)
    ToOrgPage(c *gin.Context)
}
