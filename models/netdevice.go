package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

type NetDevice struct {
	Id                  uint64              `json:"id"`
	Vendor              string              `json:"vendor"`
	ProductName         string              `json:"product_name"`
	SN                  string              `json:"sn"`
	AssetNumber         string              `json:"asset_number"`
	MgtIP               string              `json:"mgt_ip"`
	MgtUserName         string              `json:"mgt_user_name"`
	MgtPwd              string              `json:"mgt_pwd"`
	OwnerId             uint64              `json:"owner_id"`
	User                UserNoPass          `json:"user" gorm:"foreignKey:OwnerId"`
	HostName            string              `json:"host_name"`
	TagId               uint64              `json:"tag_id"`
	Tag                 Tag                 `json:"tag"`
	Location            string              `json:"location"`
	InterfaceTopologies []InterfaceTopology `json:"interface_topologies" gorm:"FOREIGNKEY:SrcNetDeviceId;ASSOCIATION_FOREIGNKEY:Id"`
	Description         string              `json:"description"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type InterfaceTopology struct {
	Id                  uint64    `json:"id"`
	SrcNetDeviceId      uint64    `json:"src_netdevice_id"`
	SrcInterfaceName    string    `json:"src_interface_name"`
	TargetNetDeviceId   uint64    `json:"target_netdevice_id"`
	TargetInterfaceName string    `json:"target_interface_name"`
	TargetNetDevice     NetDevice `json:"target_netdevice" gorm:"foreignKey:TargetNetDeviceId"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type NetDevicesRsp struct {
	Code       int          `json:"code"`
	Msg        string       `json:"message"`
	NetDevices []*NetDevice `json:"data"`
	Page       int          `json:"page"`
	PageSize   int          `json:"size"`
	Total      uint64       `json:"total"`
}

func (netDevice *NetDevice) Save() (err error) {
	if err = Db.Create(netDevice).Error; err != nil {
		return err
	}
	return
}

func (netDevice *NetDevice) GetById(id uint64) (*NetDevice, error) {
	if err := Db.Preload("Tag").Preload("User", func(db *gorm.DB) *gorm.DB { return db.Table("users") }).Where("id=?", id).First(netDevice).Error; err != nil {
		return nil, err
	}
	return netDevice, nil
}

func (netDevice *NetDevice) GetList(page int, pageSize int) (netDeviceList []*NetDevice, total uint64, err error) {
	if err = Db.Table("net_devices").Count(&total).Scopes(Paginate(page, pageSize)).Preload("InterfaceTopologies.TargetNetDevice").Preload("Tag").Preload("User", func(db *gorm.DB) *gorm.DB { return db.Table("users") }).Find(&netDeviceList).Error; err != nil {
		return nil, total, err
	}
	return netDeviceList, total, nil
}

func (netDevice *NetDevice) Delete(id uint64) (err error) {
	err = Db.Where("id=?", id).Delete(netDevice).Error
	return
}

func (netDevice *NetDevice) Update() (err error) {
	var data NetDevice
	if Db.First(&data, netDevice.Id).RecordNotFound() {
		return errors.New("resource not found")
	}
	err = Db.Save(netDevice).Error
	return
}

func (interfaceTopology *InterfaceTopology) Save() (err error) {
	if err = Db.Create(interfaceTopology).Error; err != nil {
		return err
	}
	return
}

func (interfaceTopology *InterfaceTopology) Update() (err error) {
	var data InterfaceTopology
	if Db.First(&data, interfaceTopology.Id).RecordNotFound() {
		return errors.New("resource not found")
	}
	err = Db.Save(interfaceTopology).Error
	return
}

func (interfaceTopology *InterfaceTopology) Delete(id uint64) (err error) {
	err = Db.Where("id=?", id).Delete(interfaceTopology).Error
	return
}

func (interfaceTopology *InterfaceTopology) GetList(id uint64) (interfaceTopologyList []InterfaceTopology, err error) {
	if err := Db.Table("interface_topologies").Where("src_net_device_id=?", id).Find(&interfaceTopologyList).Error; err != nil {
		return nil, err
	}
	return
}
