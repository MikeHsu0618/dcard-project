package models

import (
	orm "dcard-project/database"
	"fmt"
)

// Input 本來就沒有質, 所以不用傳地址
// create
//models.User{Name: "1111"}.Insert()

// list User{} 是值, 沒有地址
//user := models.User{}
//users, _ := user.Users()
//fmt.Printf("結果是 = %v", users)

// update
//user := models.User{Name: "太神拉"}
//users, _ := user.Update(666)
//fmt.Printf("結果是 = %v", users)

// Destroy
//user := models.User{}
//result, _ := user.Destroy(666)
//fmt.Printf("結果是 = %v", result)

type User struct {
	ID       int64  `json:"id"`        // 列名为 `id`
	OrgUrl   string `json:"org_url"`   // 列名为 `org_url`
	ShortUrl string `json:"short_url"` // 列名为 `short_url`
}

var Users []User

//添加
func (user User) Insert() (id int64, err error) {
	fmt.Printf("這是 ID 值 %v", user)

	//添加数据
	result := orm.Db.Create(&user)

	id = user.ID
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}

//取得

//列表
func (user *User) Users() (users []User, err error) {
	if err = orm.Db.First(&users).Error; err != nil {
		return
	}
	return
}

//修改
func (user *User) Update(id int64) (updateUser User, err error) {

	if err = orm.Db.Select([]string{"id", "org_url", "short_url"}).First(&updateUser, id).Error; err != nil {
		return
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Db.Model(&updateUser).Updates(&user).Error; err != nil {
		return
	}
	return
}

// Destroy 删除数据
func (user *User) Destroy(id int64) (Result User, err error) {

	if err = orm.Db.Select([]string{"id"}).First(&user, id).Error; err != nil {
		return
	}

	if err = orm.Db.Delete(&user).Error; err != nil {
		return
	}
	Result = *user
	return
}
