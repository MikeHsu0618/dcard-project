package main

import (
	"bytes"
	"dcard-project/model"
	"dcard-project/pkg/decimalconv"
	"dcard-project/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_CreateUrl(t *testing.T) {
	var w *httptest.ResponseRecorder
	var validUrl = "https://www.google.com"
	var invalidUrl = "https://www.google134f.com"

	initDBTable()
	r.POST("/", svc.CreateUrl)

	// 1. 輸入合法網址
	var jsonStr = []byte(`{"org_url":"` + validUrl + `"}`)
	w = httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var url = &model.Url{}
	var count int64
	db.Model(url).Count(&count)
	db.First(url)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 1, int(count))
	assert.Equal(t, validUrl, url.OrgUrl)

	// 2. 輸入重複網址
	jsonStr = []byte(`{"org_url":"` + url.OrgUrl + `"}`)
	w = httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	db.Model(url).Count(&count)
	assert.Equal(t, 1, int(count))
	assert.Equal(t, http.StatusOK, w.Code)

	// 3. 輸入非法網址
	jsonStr = []byte(`{"org_url":"` + invalidUrl + `"}`)
	w = httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "invalid url")
}

func Test_ToOrgPage(t *testing.T) {
	var w *httptest.ResponseRecorder
	var invalidShortUrl = "123123"

	initDBTable()
	r = gin.Default()
	r.Delims("{{{", "}}}")
	r.LoadHTMLGlob("../resources/views/*")
	r.Static("/asset", "./resources/asset")
	r.GET("/:shortUrl", svc.ToOrgPage)

	// 1. 無效的地址 -> 404 not found
	w = httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req := httptest.NewRequest(http.MethodGet, `/`+invalidShortUrl, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "無效的地址")

	// 2. 有效的地址 -> 302 轉址
	url := &model.Url{OrgUrl: "https://www.google.com"}
	db.Create(url)
	shortUrl := decimalconv.Encode(service.BasicAmount + url.ID)
	w = httptest.NewRecorder() // 取得 ResponseRecorder 物件
	req = httptest.NewRequest(http.MethodGet, `/`+shortUrl, nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusFound, w.Code)
}

// define mock type
type UrlServiceMock struct {
	mock.Mock
}

func (m *UrlServiceMock) CreateUrl(c *gin.Context) {
	m.Called(c)
}

func (m *UrlServiceMock) ToOrgPage(c *gin.Context) {
	m.Called(c)
}

func Test_CreateUrl_Mock(t *testing.T) {
	m := new(UrlServiceMock)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	m.On("CreateUrl", c)
	m.CreateUrl(c)
	m.AssertExpectations(t)
}

func Test_ToOrgPage_Mock(t *testing.T) {
	m := new(UrlServiceMock)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	m.On("ToOrgPage", c)
	m.ToOrgPage(c)
	m.AssertExpectations(t)
}
