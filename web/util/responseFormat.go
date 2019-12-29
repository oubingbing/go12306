package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	ErrorCode int `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Data interface{} `json:"data"`
	AuthorEmail string `json:"contact_email"`
}

func (r *Response) ResponseError(ctx *gin.Context)  {
	Error(fmt.Sprintf(r.ErrorMessage))
	ctx.JSON(r.ErrorCode,r)
}

func (r *Response) ResponseSuccess(ctx *gin.Context)  {
	ctx.JSON(r.ErrorCode,r)
}

func ResponseJson(ctx *gin.Context,code int,message string,data interface{})  {
	var res Response
	res.ErrorCode = code
	res.ErrorMessage = message
	res.Data = data
	res.AuthorEmail = "875307054@qq.com"

	if code != http.StatusOK {
		res.ResponseError(ctx)
	}else{
		res.ResponseSuccess(ctx)
	}
}
