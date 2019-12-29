package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"newbug/model"
	"newbug/service"
	"newbug/util"
)

func Answer(ctx *gin.Context)  {

}

func SavedeviceId(ctx *gin.Context)  {
	token := ctx.PostForm("token")
	var device model.DeviceInfo
	device.Token = token
	result,err := service.Store(&device)
	if err != nil {
		util.Error(fmt.Sprintf("注册失败：%v\n",err.Error()))
	}

	util.ResponseJson(ctx,http.StatusOK,"操作成功",result)
	return
}