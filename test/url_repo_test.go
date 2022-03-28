package main

import (
    "dcard-project/model"
    "github.com/stretchr/testify/mock"
    "testing"
)

// define mock type
type UrlRepoMock struct {
    mock.Mock
}

func (m *UrlRepoMock) Create(url *model.Url) (err error) {
    args := m.Called(url)
    return args.Error(0)
}

func (m *UrlRepoMock) GetById(urlId int64, url *model.Url) (err error) {
    args := m.Called(urlId, url)
    return args.Error(0)
}

func (m *UrlRepoMock) GetByOrgUrl(orgUrl string) (url *model.Url, err error) {
    args := m.Called(orgUrl)
    return args.Get(0).(*model.Url), args.Error(1)
}

func (m *UrlRepoMock) Lock(key string) bool {
    args := m.Called(key)
    return args.Bool(0)
}

func (m *UrlRepoMock) UnLock(key string) int64 {
    args := m.Called(key)

    return int64(args.Int(0))
}

func (m *UrlRepoMock) GetCache(key int64) (result string, err error) {
    args := m.Called(key)
    return args.String(0), args.Error(1)
}

func (m *UrlRepoMock) SetCache(key int64, orgUrl string) {
    m.Called(key, orgUrl)
}

func Test_Create(t *testing.T) {
    m := new(UrlRepoMock)

    url := &model.Url{
        ID:     0,
        OrgUrl: `https:\\www.google.com`,
    }

    m.On("Create", url).Return(nil)
    m.Create(url)

    m.AssertExpectations(t)
}

func Test_GetById(t *testing.T) {
    m := new(UrlRepoMock)
    url := &model.Url{}
    m.On("GetById", int64(1), url).Return(nil)
    m.GetById(1, url)

    m.AssertExpectations(t)
}

func Test_GetByOrgUrl(t *testing.T) {
    m := new(UrlRepoMock)
    orgUrl := "https://www.google.com"
    url := &model.Url{}
    m.On("GetByOrgUrl", orgUrl).Return(url, nil)
    m.GetByOrgUrl(orgUrl)

    m.AssertExpectations(t)
}

func Test_Lock(t *testing.T) {
    m := new(UrlRepoMock)
    key := "key"
    m.On("Lock", key).Return(true)
    m.Lock(key)

    m.AssertExpectations(t)
}

func Test_UnLock(t *testing.T) {
    m := new(UrlRepoMock)
    key := "key"
    m.On("UnLock", key).Return(1)
    m.UnLock(key)

    m.AssertExpectations(t)
}

func Test_GetCache(t *testing.T) {
    m := new(UrlRepoMock)
    key := int64(1)
    result := "https://www.google.com"
    m.On("GetCache", key).Return(result, nil)
    m.GetCache(key)

    m.AssertExpectations(t)
}

func Test_SetCache(t *testing.T) {
    m := new(UrlRepoMock)
    key := int64(1)
    orgUrl := "https://www.google.com"
    m.On("SetCache", key, orgUrl)
    m.SetCache(key, orgUrl)

    m.AssertExpectations(t)
}
