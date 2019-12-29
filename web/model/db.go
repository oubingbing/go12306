package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"newbug/util"
)

func Connect() *gorm.DB {
	configString ,err := util.GetMysqlConfig()
	fmt.Println(configString)
	if err != nil{
		util.Info(fmt.Sprintf("获取数据失败：%v\n",err.Error()))
	}

	db, err := gorm.Open("mysql", configString)
	if err != nil {
		util.Info(fmt.Sprintf("连接数据库错误：%v\n",err.Error()))
	}

	return db
}

