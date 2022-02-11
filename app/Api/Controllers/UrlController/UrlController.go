package UrlController

import (
	"context"
	"dcard-project/app/Constant"
	"dcard-project/app/Logic/DecimalConvert"
	"dcard-project/app/Logic/Log"
	. "dcard-project/database"
	"dcard-project/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var ctx = context.Background()
var url = &models.Url{}

func Show(c *gin.Context) {
	shortUrl := c.Param("shortUrl")
	Db.Where("short_url", shortUrl).First(&models.Urls)

	c.JSON(200, gin.H{
		"data":    &models.Urls,
		"status":  200,
		"message": "success",
	})
}

func Create(c *gin.Context) {
	// 接收參數
	if err := c.ShouldBind(&url); err != nil {
		c.JSON(404, gin.H{
			"data":    err.Error(),
			"message": "input error",
		})
		return
	}
	// 檢查原網址
	res, err := http.Get(url.OrgUrl)
	if err != nil || res.StatusCode != 200 {
		c.JSON(404, gin.H{
			"data":    "",
			"message": "invalid url",
		})
		return
	}
	// 產生縮網址
	shortUrl := getShortUrl()
	// 保存三十天過期
	_, err = Redis.Set(context.Background(), shortUrl, url.OrgUrl, 30*24*time.Hour).Result()
	if err != nil {
		Log.Error.Println("Redis Set Url Error", err.Error())
		return
	}
	// 檢查是否已重複
	result := Db.Create(&models.Url{
		OrgUrl:   url.OrgUrl,
		ShortUrl: shortUrl,
	})
	// 已存在則返回
	err = result.Error
	if err != nil && strings.Contains(err.Error(), "duplicate") {
		duplicateUrl := &models.Url{}
		Db.Where("org_url", url.OrgUrl).First(duplicateUrl)
		c.JSON(200, gin.H{
			"data":    duplicateUrl.ShortUrl,
			"message": "success",
		})
		return
	}
	if err != nil {
		c.JSON(404, gin.H{
			"data":    "",
			"message": "create fail",
		})
		return
	}

	c.JSON(200, gin.H{
		"data":    shortUrl,
		"message": "success",
	})
}

func ToOrgPage(c *gin.Context) {
	result, _ := Redis.Get(ctx, c.Param("shortUrl")).Result()
	// 使用快取
	if len(result) != 0 {
		println("我用快取拉 我發達了", result)
		c.Redirect(http.StatusFound, result)
		return
	}

	// 使用資料庫
	for {
		// 上鎖
		if Lock(Constant.LockKey) == false {
			time.Sleep(100 * time.Millisecond)
			println("waiting")
			continue
		}
		println("succeed aaa")
		Db.Where("short_url", c.Param("shortUrl")).First(&url)

		if len(url.OrgUrl) == 0 {
			c.HTML(
				http.StatusNotFound,
				"404.html",
				gin.H{"title": "無效的地址"},
			)
			return
		}
		UnLock(Constant.LockKey)
		break
	}
	c.Redirect(http.StatusFound, url.OrgUrl)
}

func getShortUrl() (shortUrl string) {
	index, err := url.Count()
	if err != nil {
		Log.Error.Println(err)
		return
	}
	shortUrl = DecimalConvert.Encode(Constant.BasicAmount + index)
	return
}
