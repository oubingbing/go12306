package service

import "newbug/model"

func Store(device *model.DeviceInfo) (int64,error) {
	createResult := model.Connect().Create(&device)
	return createResult.RowsAffected,createResult.Error
}