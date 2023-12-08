package models

import (
	"errors"
	"fmt"
	"time"
)

type User struct {
	Id        uint16 `gorm:"primary_key" json:"id"`
	Username  string `json:"username"`
	NameAlias string `json:"name_alias"`
	Password  string `json:"password"`
	Phone     string `json:"phone" gorm:"default: null"`
	Mail      string `json:"mail" gorm:"default: null"`
	Role      string `json:"role"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserNoPass struct {
	Id        uint16 `gorm:"primary_key" json:"id"`
	Username  string `json:"username"`
	NameAlias string `json:"name_alias"`
	Phone     string `json:"phone" gorm:"default: null"`
	Mail      string `json:"mail" gorm:"default: null"`
	Role      string `json:"role"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRsp struct {
	Code     int           `json:"code"`
	Msg      string        `json:"message"`
	Users    []*UserNoPass `json:"data"`
	Page     int           `json:"page"`
	PageSize int           `json:"size"`
	Total    uint64        `json:"total"`
}

func (user *User) CheckUser() (exist bool, userRole string, userId uint16, username string) {
	Db.Select("id,role,username").Where(User{Username: user.Username, Password: user.Password}).First(&user)

	if user.Id > 0 {
		exist = true
		userRole = user.Role
		userId = user.Id
		username = user.Username
		return
	}
	return false, "", userId, username
}

func (user *User) GetList(page int, pageSize int) (userList []*UserNoPass, total uint64, err error) {
	if err = Db.Table("users").Count(&total).Scopes(Paginate(page, pageSize)).Find(&userList).Error; err != nil {
		return nil, total, err
	}
	return userList, total, nil
}

func (user *User) Add() (err error) {
	var data User
	Db.Where(User{Username: user.Username}).First(&data)
	if data.Id > 0 {
		errMsg := fmt.Sprintf("User %s already exist.", user.Username)
		return errors.New(errMsg)
	}
	if err = Db.Create(user).Error; err != nil {
		return err
	}
	return
}

func (user *User) Update() (err error) {
	var data User
	if Db.First(&data, user.Id).RecordNotFound() {
		return errors.New("resource not found")
	}
	err = Db.Model(&user).Updates(map[string]interface{}{"name_alias": user.NameAlias, "phone": user.Phone, "mail": user.Mail}).Error
	return
}

func (user *User) Delete(id uint64) (err error) {
	err = Db.Where("id=?", id).Delete(user).Error
	return
}
