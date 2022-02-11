package models

import (
	. "dcard-project/database"
	"fmt"
)

// Input 本來就沒有質, 所以不用傳地址
// create
//models.Url{Name: "1111"}.Insert()

// list Url{} 是值, 沒有地址
//url := models.Url{}
//urls, _ := url.Urls()
//fmt.Printf("結果是 = %v", urls)

// update
//url := models.Url{Name: "太神拉"}
//urls, _ := url.Update(666)
//fmt.Printf("結果是 = %v", urls)

// Destroy
//url := models.Url{}
//result, _ := url.Destroy(666)
//fmt.Printf("結果是 = %v", result)

type Url struct {
	ID       int64  `json:"id" form:"id"`                              // 列名为 `id`
	OrgUrl   string `json:"org_url" form:"org_url" binding:"required"` // 列名为 `org_url`
	ShortUrl string `json:"short_url"`                                 // 列名为 `short_url`
}

type ApiUrl struct {
	ID     int64  `json:"id"`      // 列名为 `id`
	OrgUrl string `json:"org_url"` // 列名为 `org_url`
}

var Urls []Url

//添加
func (url Url) Insert() (id int64, err error) {
	fmt.Printf("這是 ID 值 %v", url)

	//添加数据
	result := Db.Create(&url)

	id = url.ID
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

//取得

//列表
func (url *Url) Urls() (urls []Url, err error) {
	if err = Db.First(&urls).Error; err != nil {
		return
	}
	return
}

func (url *Url) Count() (count int64, err error) {
	if err = Db.Model(&url).Count(&count).Error; err != nil {
		return
	}
	return
}

//修改
func (url *Url) Update(id int64) (updateUrl Url, err error) {

	if err = Db.Select([]string{"id", "org_url", "short_url"}).First(&updateUrl, id).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = Db.Model(&updateUrl).Updates(&url).Error; err != nil {
		return
	}
	return
}

// Destroy 删除数据
func (url *Url) Destroy(id int64) (Result Url, err error) {

	if err = Db.Select([]string{"id"}).First(&url, id).Error; err != nil {
		return
	}

	if err = Db.Delete(&url).Error; err != nil {
		return
	}
	Result = *url
	return
}
