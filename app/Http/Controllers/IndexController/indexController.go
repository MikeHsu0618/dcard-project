package IndexController

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Show(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{"title": "縮網址產生器"},
	)
}
