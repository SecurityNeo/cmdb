package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type Host struct {
	Id              uint64     `json:"id"`
	HdmIP           string     `json:"hdm_ip"`
	HdmUserName     string     `json:"hdm_user_name"`
	HdmPwd          string     `json:"hdm_password"`
	Vendor          string     `json:"vendor"`
	ProductName     string     `json:"product_name"`
	Arch            string     `json:"arch"`
	SN              string     `json:"sn"`
	AssetNumber     string     `json:"asset_number"`
	HardwareCPU     int        `json:"hardware_cpu"`
	HardwareMem     int        `json:"hardware_mem"`
	HardwareStorage string     `json:"hardware_storage"`
	MgtIP           string     `json:"mgt_ip"`
	MgtRootPwd      string     `json:"mgt_root_pwd"`
	StorageIP       string     `json:"storage_ip"`
	OwnerId         uint64     `json:"owner_id"`
	User            UserNoPass `json:"user" gorm:"foreignKey:OwnerId"`
	HostName        string     `json:"host_name"`
	TagId           uint64     `json:"tag_id"`
	Tag             Tag        `json:"tag"`
	Location        string     `json:"location"`
	Description     string     `json:"description"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type HostsRsp struct {
	Code     int     `json:"code"`
	Msg      string  `json:"message"`
	Hosts    []*Host `json:"data"`
	Page     int     `json:"page"`
	PageSize int     `json:"size"`
	Total    uint64  `json:"total"`
}

func (host *Host) Save() (err error) {
	if err = Db.Create(host).Error; err != nil {
		return err
	}
	return
}

func (host *Host) GetById(id uint64) (*Host, error) {
	if err := Db.Preload("Tag").Preload("User", func(db *gorm.DB) *gorm.DB { return db.Table("users") }).Where("id=?", id).First(host).Error; err != nil {
		return nil, err
	}
	return host, nil
}

func (host *Host) GetList(page int, pageSize int, TagIdStr string) (hostList []*Host, total uint64, err error) {

	if TagIdStr == "" {
		if err = Db.Table("hosts").Count(&total).Scopes(Paginate(page, pageSize)).Preload("Tag").Preload("User", func(db *gorm.DB) *gorm.DB { return db.Table("users") }).Find(&hostList).Error; err != nil {
			return nil, total, err
		}
	} else {
		tagId, err := strconv.ParseUint(TagIdStr, 10, 64)
		if err != nil {
			return nil, total, err
		}
		if err = Db.Table("hosts").Count(&total).Scopes(Paginate(page, pageSize)).Preload("Tag").Preload("User", func(db *gorm.DB) *gorm.DB { return db.Table("users") }).Where("tag_id=?", tagId).Find(&hostList).Error; err != nil {
			return nil, total, err
		}
	}

	return hostList, total, nil
}

func (host *Host) Delete(id uint64) (err error) {
	err = Db.Where("id=?", id).Unscoped().Delete(host).Error
	return
}

func (host *Host) Update() (err error) {
	var data Host
	if Db.First(&data, host.Id).RecordNotFound() {
		return errors.New("resource not found")
	}
	err = Db.Save(host).Error
	return
}
