package models

import (
	"errors"
	"time"
)

type Tag struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TagRsp struct {
	Code     int    `json:"code"`
	Msg      string `json:"message"`
	Tags     []*Tag `json:"data"`
	Page     int    `json:"page"`
	PageSize int    `json:"size"`
	Total    uint64 `json:"total"`
}

func (tag *Tag) Save() (err error) {
	if err = Db.Create(tag).Error; err != nil {
		return err
	}
	return
}

func (tag *Tag) GetById(id uint64) (*Tag, error) {
	if err := Db.Where("id=?", id).First(&tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

func (tag *Tag) GetList(page int, pageSize int) (tagList []*Tag, total uint64, err error) {
	if err = Db.Table("tags").Count(&total).Scopes(Paginate(page, pageSize)).Find(&tagList).Error; err != nil {
		return nil, total, err
	}
	return tagList, total, nil
}

func (tag *Tag) Delete(id uint64) (err error) {
	err = Db.Where("id=?", id).Delete(tag).Error
	return
}

func (tag *Tag) Update() (err error) {
	var data Tag
	if Db.First(&data, tag.Id).RecordNotFound() {
		return errors.New("resource not found")
	}
	err = Db.Save(tag).Error
	return
}
