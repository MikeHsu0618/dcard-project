package v1

import (
    "context"
    . "dcard-project/database"
    _ "dcard-project/docs"
    "dcard-project/models"
    "dcard-project/pkg/decimalconv"
    "dcard-project/pkg/goquery"
    "dcard-project/pkg/httputil"
    "github.com/gin-gonic/gin"
    "net/http"
    "strconv"
    "strings"
    "time"
)

const (
    LockKey     = "lock_key"
    BasicAmount = int64(20000)
)

var ctx = context.Background()

// @Summary 產生短網址
// @Description 請輸入合法原網址
// @Tags Url
// @Accept json
// @Produce json
// @Param url body models.createUrl true "Get Short Url"
// @Success 200 {object} models.ApiUrl
// @Router / [post]
func CreateUrl(c *gin.Context) {
    var url = &models.Url{}
    var shortUrl string
    // 接收參數
    if err := c.ShouldBind(&url); err != nil {
        httputil.NewError(c, 404, err.Error())
        return
    }
    // 檢查原網址
    res, err := http.Get(url.OrgUrl)
    if err != nil || res.StatusCode != http.StatusOK {
        httputil.NewError(c, http.StatusNotFound, "invalid url")
        return
    }
    meta := goquery.GetHtmlMeta(res.Body)
    // 檢查是否已重複
    result := Db.Create(url)
    // 已存在則返回
    if err := result.Error; err != nil && strings.Contains(err.Error(), "duplicate") {
        duplicateUrl := &models.Url{}
        Db.Where("org_url", url.OrgUrl).First(duplicateUrl)
        shortUrl = decimalconv.Encode(BasicAmount + duplicateUrl.ID)
        data := models.ApiUrl{
            ShortUrl: shortUrl,
            Meta:     meta,
        }
        httputil.NewSuccess(c, data)
        return
    }
    if result.Error != nil {
        httputil.NewError(c, http.StatusNotFound, "create fail")
        return
    }

    //產生縮網址
    shortUrl = decimalconv.Encode(BasicAmount + url.ID)
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

    httputil.NewSuccess(c, data)
}

func ToOrgPage(c *gin.Context) {
    var url = &models.Url{}
    index := decimalconv.Decode(c.Param("shortUrl")) - BasicAmount
    // 使用快取
    if result, _ := Redis.Get(ctx, strconv.FormatInt(index, 10)).Result(); len(result) != 0 {
        c.Redirect(http.StatusFound, result)
        return
    }

    // 使用資料庫
    for {
        // 上鎖
        if Lock(LockKey) == false {
            time.Sleep(100 * time.Millisecond)
            continue
        }
        Db.Where("id", index).First(&url)
        if len(url.OrgUrl) == 0 {
            UnLock(LockKey)
            c.HTML(
                http.StatusNotFound,
                "404.html",
                gin.H{"title": "無效的地址"},
            )
            return
        }
        Redis.Set(context.Background(), strconv.FormatInt(index, 10), url.OrgUrl, 30*24*time.Hour)
        UnLock(LockKey)
        break
    }

    // 保存三十天過期
    Redis.Set(context.Background(), strconv.FormatInt(index, 10), url.OrgUrl, 30*24*time.Hour)
    c.Redirect(http.StatusFound, url.OrgUrl)
}
