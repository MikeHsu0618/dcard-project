package Middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		// 可将将* 替换为指定的域名
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}

//func Cors() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		method := c.Request.Method
//		//origin := "*" //請求頭部
//		//fmt.Printf("%v", origin)
//		////接收客戶端傳送的origin （重要！）
//		//c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
//		////伺服器支援的所有跨域請求的方法
//		//c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
//		////允許跨域設定可以返回其他子段，可以自定義欄位
//		//c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
//		//// 允許瀏覽器（客戶端）可以解析的頭部 （重要）
//		//c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
//		////設定快取時間
//		//c.Header("Access-Control-Max-Age", "172800")
//		////允許客戶端傳遞校驗資訊比如 cookie (重要)
//		//c.Header("Access-Control-Allow-Credentials", "true")
//
//		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
//		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
//		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
//		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
//		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
//
//		//允許型別校驗
//		if method == "OPTIONS" {
//			c.JSON(http.StatusOK, "ok!")
//		}
//
//		defer func() {
//			if err := recover(); err != nil {
//				log.Printf("Panic info is: %v", err)
//			}
//		}()
//
//		c.Next()
//	}
//}
