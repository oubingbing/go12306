package model

import "github.com/jinzhu/gorm"

type DeviceInfo struct {
	gorm.Model
	Token string `form:"token"`
}

func (DeviceInfo) TableName() string  {
	return "device_info"
}