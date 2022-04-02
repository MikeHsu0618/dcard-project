package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"dcard-project/pkg/httputil"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const (
	IPLimitPeriod           = 60
	IPLimitTimeFormat       = "2006-01-02 15:04:05"
	IPLimitMaximum    int64 = 100
)

var ctx = context.Background()
var key string
var amount int64
var rdsClient *redis.Client

func Setup(client *redis.Client) {
	rdsClient = client
}

func IPLimitIntercept() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now().Unix()
		amount = getCurrentAmount(c)
		rdsClient.Set(ctx, key, amount+1, IPLimitPeriod*time.Second)
		reset := time.Unix(now+IPLimitPeriod, 0).Format(IPLimitTimeFormat)
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(IPLimitMaximum-amount, 10))
		c.Header("X-RateLimit-Reset", reset)
	}
}

func getCurrentAmount(c *gin.Context) (amount int64) {
	var err error
	key = c.Request.Method + "-" + c.ClientIP()
	amount, err = rdsClient.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0
	}
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if amount >= IPLimitMaximum {
		httputil.NewError(c, http.StatusInternalServerError, err.Error())
		return
	}

	return
}
