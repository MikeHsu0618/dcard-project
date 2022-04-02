package service

import (
	"net/http"
	"strings"
	"time"

	"dcard-project/internal/repository"
	"dcard-project/pkg/decimalconv"
	"dcard-project/pkg/goquery"
	"dcard-project/pkg/httputil"
	"github.com/gin-gonic/gin"
)

const (
	LockKey       = "lock_key"
	BasicAmount   = int64(20000)
	UserAgentInfo = `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.83 Safari/537.36`
)

type UrlService struct {
	repo repository.UrlRepository
}

func NewUrlService(urlRepo repository.UrlRepository) repository.UrlService {
	return &UrlService{repo: urlRepo}
}

// @Summary 產生短網址
// @Description 請輸入合法原網址
// @Tags Url
// @Accept json
// @Produce json
// @Param url body repository.CreateUrl true "Get Short Url"
// @Success 200 {object} repository.ApiUrl
// @Router / [post]
func (s *UrlService) CreateUrl(c *gin.Context) {
	var url = &repository.Url{}
	var shortUrl string
	// 接收參數
	if err := c.ShouldBind(&url); err != nil {
		httputil.NewError(c, 404, err.Error())
		return
	}
	// 檢查原網址
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, url.OrgUrl, nil)
	req.Header.Set("User-Agent", UserAgentInfo)
	res, err := client.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		httputil.NewError(c, http.StatusBadRequest, "invalid url")
		return
	}
	meta := goquery.GetHtmlMeta(res.Body)
	// 檢查是否已存在
	err = s.repo.Create(url)
	if err != nil &&
		(strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "UNIQUE")) {
		duplicateUrl, err := s.repo.GetByOrgUrl(url.OrgUrl)
		if err != nil {
			httputil.NewError(c, http.StatusServiceUnavailable, "data error")
			return
		}
		shortUrl = decimalconv.Encode(BasicAmount + duplicateUrl.ID)
		httputil.NewSuccess(c, repository.ApiUrl{
			ShortUrl: shortUrl,
			Meta:     meta,
		})
		return
	} else if err != nil {
		httputil.NewError(c, http.StatusServiceUnavailable, "create fail")
		return
	}

	// 產生短網址
	shortUrl = decimalconv.Encode(BasicAmount + url.ID)
	data := repository.ApiUrl{
		ShortUrl: shortUrl,
		Meta:     meta,
	}
	// 保存三十天過期
	s.repo.SetCache(url.ID, url.OrgUrl)

	httputil.NewSuccess(c, data)
}

func (s *UrlService) ToOrgPage(c *gin.Context) {
	var url = &repository.Url{}
	index := decimalconv.Decode(c.Param("shortUrl")) - BasicAmount
	// 使用快取
	if orgUrl, err := s.repo.GetCache(index); err == nil && len(orgUrl) != 0 {
		// 保存三十天過期
		s.repo.SetCache(index, orgUrl)
		c.Redirect(http.StatusFound, orgUrl)
		return
	}

	// 使用資料庫
	for {
		// 上鎖
		if s.repo.Lock(LockKey) == false {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if err := s.repo.GetById(index, url); err != nil || len(url.OrgUrl) == 0 {
			s.repo.UnLock(LockKey)
			c.HTML(
				http.StatusNotFound,
				"404.html",
				gin.H{"title": "無效的地址"},
			)
			return
		}
		s.repo.SetCache(index, url.OrgUrl)
		s.repo.UnLock(LockKey)
		break
	}
	c.Redirect(http.StatusFound, url.OrgUrl)
}
