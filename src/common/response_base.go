package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//定义一个通用的返回结构体

type BaseResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"result"`
	Type    string      `json:"type"`
}

const (
	Success = 0
	Error   = 7
)

func Result(code int, data interface{}, message string, c *gin.Context) {

	c.JSON(http.StatusOK, BaseResp{
		Code:    code,
		Message: message,
		Data:    data,
		Type:    "",
	})
}
func Ok(c *gin.Context) {
	Result(Success, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(Success, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(Success, data, "查询成功", c)
}

func Fail(c *gin.Context) {
	Result(Error, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(Error, map[string]interface{}{}, message, c)
}
func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(Error, data, message, c)
}
func Result4xx(code int, data interface{}, message string, c *gin.Context) {
	c.JSON(http.StatusBadRequest, BaseResp{
		Code:    code,
		Message: message,
		Data:    data,
		Type:    "",
	})
}

func Result401(data interface{}, message string, c *gin.Context) {
	Result4xx(401, data, message, c)
}

func ReBadFailWithMessage(message string, c *gin.Context) {
	Result4xx(Error, map[string]interface{}{}, message, c)
}

func ReBadFailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result4xx(Error, data, message, c)
}

func Re401FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result401(data, message, c)
}
