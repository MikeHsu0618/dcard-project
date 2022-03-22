package models

import (
    "dcard-project/pkg/goquery"
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

var Urls []Url
