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

func Create(c *gin.Context) {
	var url = &models.Url{}

	// 接收參數
	if err := c.ShouldBind(&url); err != nil {
		c.JSON(404, gin.H{
			"data":    err.Error(),
			"message": "input error",
		})
		return
	}
	// 檢查原網址
	if res, err := http.Get(url.OrgUrl); err != nil || res.StatusCode != 200 {
		c.JSON(404, gin.H{
			"data":    "",
			"message": "invalid url",
		})
		return
	}
	// 產生縮網址
	shortUrl := getShortUrl()
	// 保存三十天過期
	if _, err := Redis.Set(
		context.Background(),
		shortUrl,
		url.OrgUrl,
		30*24*time.Hour).Result(); err != nil {
		Log.Error.Println("Redis Set Url Error", err.Error())
		return
	}
	// 檢查是否已重複
	result := Db.Create(&models.Url{
		OrgUrl:   url.OrgUrl,
		ShortUrl: shortUrl,
	})
	// 已存在則返回
	if err := result.Error; err != nil && strings.Contains(err.Error(), "duplicate") {
		duplicateUrl := &models.Url{}
		Db.Where("org_url", url.OrgUrl).First(duplicateUrl)
		c.JSON(200, gin.H{
			"data":    duplicateUrl.ShortUrl,
			"message": "success",
		})
		return
	}
	if result.Error != nil {
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
	var url = &models.Url{}
	// 使用快取
	if result, _ := Redis.Get(ctx, c.Param("shortUrl")).Result(); len(result) != 0 {
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
		println("success getting data")
		Db.Where("short_url", c.Param("shortUrl")).First(&url)
		if len(url.OrgUrl) == 0 {
			c.HTML(
				http.StatusNotFound,
				"404.html",
				gin.H{"title": "無效的地址"},
			)
			UnLock(Constant.LockKey)
			return
		}

		UnLock(Constant.LockKey)
		break
	}

	// 保存三十天過期
	Redis.Set(context.Background(), url.ShortUrl, url.OrgUrl, 30*24*time.Hour)
	c.Redirect(http.StatusFound, url.OrgUrl)
}

func getShortUrl() (shortUrl string) {
	var url = &models.Url{}
	index, err := url.Count()
	if err != nil {
		Log.Error.Println(err)
		return
	}
	shortUrl = DecimalConvert.Encode(Constant.BasicAmount + index)
	return
}
