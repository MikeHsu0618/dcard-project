package UrlControllers

import (
	"dcard-project/app/Logic/decimalConvert"
	. "dcard-project/database"
	"dcard-project/models"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

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
	url := models.Url{}
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

	// TODO 取得 autoIncrement index id
	//var idInt int64 = 1
	// 產生亂數
	rand.Seed(time.Now().Unix())
	randNum := int64(rand.Intn(10000000))
	var basicAmount int64 = int64(20000)
	shortUrl := decimalConvert.Encode(basicAmount + randNum)

	// 檢查是否已重複
	result := Db.Create(&models.Url{
		OrgUrl:   url.OrgUrl,
		ShortUrl: shortUrl,
	})
	if err := result.Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			duplicateUrl := &models.Url{}
			Db.Where("org_url", url.OrgUrl).First(duplicateUrl)
			c.JSON(200, gin.H{
				"data":    duplicateUrl.ShortUrl,
				"message": "success",
			})
			return
		}

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
	url := &models.Url{}
	Db.Where("short_url", c.Param("shortUrl")).First(url)
	c.Redirect(http.StatusFound, url.OrgUrl)
}
