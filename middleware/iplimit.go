package middleware

import (
    "context"
    "dcard-project/pkg/httputil"
    . "dcard-project/pkg/redis"
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
        client := GetRedisClient()
        amount = getCurrentAmount(c)
        client.Set(ctx, key, amount+1, IPLimitPeriod*time.Second)
        reset := time.Unix(now+IPLimitPeriod, 0).Format(IPLimitTimeFormat)
        c.Header("X-RateLimit-Remaining", strconv.FormatInt(IPLimitMaximum-amount, 10))
        c.Header("X-RateLimit-Reset", reset)
    }
}

func getCurrentAmount(c *gin.Context) (amount int64) {
    var err error
    client := GetRedisClient()
    //key = c.Request.URL.Path + "-" + c.Request.Method + "-" + c.ClientIP()
    key = c.Request.Method + "-" + c.ClientIP()
    amount, err = client.Get(ctx, key).Int64()
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
