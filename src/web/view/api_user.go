package view

import (
	"bigyunwei-backend/src/common"
	"bigyunwei-backend/src/models"
	"errors"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

func UserLogin(c *gin.Context) {
	var user models.UserLoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		common.FailWithMessage(err.Error(), c)
		return
	}
	//校验参数
	if err := validate.Struct(user); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			common.FailWithDetailed(validationErrors.Translate(trans), "请求出错", c)
			return
		}
		common.FailWithMessage(err.Error(), c)
		return
	}

	//校验用户
	dbUser := &models.User{
		Username: user.Username,
		Password: user.Password,
	}

	//生成jwt的token
	models.TokenNext(dbUser, c)

}
