package Middleware

import (
	"context"
	. "dcard-project/database"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strconv"
	"time"
)

const (
	IPLimitPeriod           = 3600
	IPLimitTimeFormat       = "2006-01-02 15:04:05"
	IPLimitMaximum    int64 = 1000
)

var ctx = context.Background()
var key string
var amount int64

func IPLimitIntercept() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now().Unix()
		amount = getCurrentAmount(c)
		Redis.Set(ctx, key, amount+1, IPLimitPeriod*time.Second)
		reset := time.Unix(now+IPLimitPeriod, 0).Format(IPLimitTimeFormat)
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(IPLimitMaximum-amount, 10))
		c.Header("X-RateLimit-Reset", reset)
	}
}

func getCurrentAmount(c *gin.Context) (amount int64) {
	key = c.Request.URL.Path + "-" + c.Request.Method + "-" + c.ClientIP()

	var err error
	amount, err = Redis.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0
	}
	if err != nil {
		sendResponse(c, http.StatusInternalServerError, err)
		return
	}
	if amount >= IPLimitMaximum {
		sendResponse(c, http.StatusTooManyRequests, err)
		return
	}

	return
}

func sendResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, data)
	c.Abort()
}
