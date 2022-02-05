package UrlControllers

import (
	. "dcard-project/database"
	"dcard-project/models"
	"github.com/gin-gonic/gin"
	"net/http"
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
	url := models.Url{}
	if err := c.ShouldBind(&url); err != nil {
		c.JSON(404, gin.H{
			"data":    err.Error(),
			"message": "input error",
		})
		return
	}

	result := Db.Create(&models.Url{
		OrgUrl:   url.OrgUrl,
		ShortUrl: url.OrgUrl,
	})
	if err := result.Error; err != nil {
		c.JSON(404, gin.H{
			"data":    err.Error(),
			"message": "create fail",
		})
		return
	}

	c.JSON(200, gin.H{
		"data":    url.OrgUrl,
		"message": "success",
	})
}

func ToOrgPage(c *gin.Context) {
	url := &models.Url{}
	Db.Where("short_url", c.Param("shortUrl")).First(url)
	c.Redirect(http.StatusFound, url.OrgUrl)
}
