package UrlController

import (
	"context"
	"dcard-project/app/Constant"
	"dcard-project/app/Logic/DecimalConvert"
	"dcard-project/app/Logic/Goquery"
	. "dcard-project/database"
	"dcard-project/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var ctx = context.Background()

func Create(c *gin.Context) {
	var url = &models.Url{}
	var shortUrl string
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
	meta := Goquery.GetHtmlMeta(res.Body)

	// 檢查是否已重複
	result := Db.Create(url)
	// 已存在則返回
	if err := result.Error; err != nil && strings.Contains(err.Error(), "duplicate") {
		duplicateUrl := &models.Url{}
		Db.Where("org_url", url.OrgUrl).First(duplicateUrl)
		shortUrl = DecimalConvert.Encode(Constant.BasicAmount + duplicateUrl.ID)
		data := models.ApiUrl{
			ShortUrl: shortUrl,
			Meta:     meta,
		}
		c.JSON(200, gin.H{
			"data":    data,
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

	//產生縮網址
	shortUrl = DecimalConvert.Encode(Constant.BasicAmount + url.ID)
	data := models.ApiUrl{
		ShortUrl: shortUrl,
		Meta:     meta,
	}
	//保存三十天過期
	if _, err := Redis.Set(
		context.Background(),
		strconv.FormatInt(url.ID, 10),
		url.OrgUrl,
		30*24*time.Hour).Result(); err != nil {
		println("Redis Set Url Error", err.Error())
		return
	}

	c.JSON(200, gin.H{
		"data":    data,
		"message": "success",
	})
}

func ToOrgPage(c *gin.Context) {
	var url = &models.Url{}
	index := DecimalConvert.Decode(c.Param("shortUrl")) - Constant.BasicAmount
	// 使用快取
	if result, _ := Redis.Get(ctx, strconv.FormatInt(index, 10)).Result(); len(result) != 0 {
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
		Db.Where("id", index).First(&url)
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
	Redis.Set(context.Background(), strconv.FormatInt(index, 10), url.OrgUrl, 30*24*time.Hour)
	c.Redirect(http.StatusFound, url.OrgUrl)
}
