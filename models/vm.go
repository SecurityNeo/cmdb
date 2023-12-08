package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

type Vm struct {
	Id          uint64     `json:"id"`
	MgtIP       string     `json:"mgt_ip"`
	MgtRootPwd  string     `json:"mgt_root_pwd"`
	StorageIP   string     `json:"storage_ip"`
	OwnerId     uint64     `json:"owner_id"`
	User        UserNoPass `json:"user" gorm:"foreignKey:OwnerId"`
	HostName    string     `json:"host_name"`
	TagId       uint64     `json:"tag_id"`
	Tag         Tag        `json:"tag"`
	HostId      uint64     `json:"host_id"`
	Host        Host       `json:"host"`
	Description string     `json:"description"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type VmsRsp struct {
	Code     int    `json:"code"`
	Msg      string `json:"message"`
	Vms      []*Vm  `json:"data"`
	Page     int    `json:"page"`
	PageSize int    `json:"size"`
	Total    uint64 `json:"total"`
}

func (tag *Vm) Save() (err error) {
	if err = Db.Create(tag).Error; err != nil {
		return err
	}
	return
}

func (tag *Vm) GetById(id uint64) (*Vm, error) {
	if err := Db.Preload("Host").Preload("Tag").Preload("User", func(db *gorm.DB) *gorm.DB { return db.Table("users") }).Where("id=?", id).First(&tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

func (tag *Vm) GetList(page int, pageSize int) (vmList []*Vm, total uint64, err error) {
	if err = Db.Table("vms").Count(&total).Scopes(Paginate(page, pageSize)).Preload("Host").Preload("Tag").Preload("User", func(db *gorm.DB) *gorm.DB { return db.Table("users") }).Find(&vmList).Error; err != nil {
		return nil, total, err
	}
	return vmList, total, nil
}

func (tag *Vm) Delete(id uint64) (err error) {
	err = Db.Where("id=?", id).Delete(tag).Error
	return
}

func (tag *Vm) Update() (err error) {
	var data Vm
	if Db.First(&data, tag.Id).RecordNotFound() {
		return errors.New("resource not found")
	}
	err = Db.Save(tag).Error
	return
}
