package main

import (
	orm "dcard-project/database"
	"dcard-project/models"
	"fmt"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int64  `json:"id"`        // 列名为 `id`
	OrgUrl   string `json:"org_url"`   // 列名为 `org_url`
	ShortUrl string `json:"short_url"` // 列名为 `short_url`
}

var Users []User

func main() {
	// list User{} 是值, 沒有地址
	user := models.User{}
	users, _ := user.Users()
	fmt.Printf("結果是 = %v", users)

	//setRouter()
}

func getUrl(c *gin.Context) {
	//shortUrl := c.Param("shortUrl")
	//list User{} 是值, 沒有地址
	user := models.User{}
	users := orm.Db.First(user)
	if users.Error != nil {
		fmt.Printf("結果是 = %v", users.Error)
		return
	}

	fmt.Printf("結果是 = %v", users)
	c.JSON(200, gin.H{"message": users})
}

func setRouter() {
	router := gin.Default()

	v1Router := router.Group("/v1/")
	{
		v1Router.GET("/:shortUrl", getUrl)
	}

	router.Run()
}
