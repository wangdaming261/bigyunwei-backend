package view

import (
	"bigyunwei-backend/src/common"
	"bigyunwei-backend/src/config"
	"bigyunwei-backend/src/models"
	"errors"
	"fmt"
	"go.uber.org/zap"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

func UserLogin(c *gin.Context) {
	var user models.UserLoginRequest
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
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

	dbUser, err := models.CheckUserPassword(&user)
	if err != nil {
		sc.Logger.Error("登录失败，校验用户失败", zap.Error(err))
		common.ReBadFailWithMessage("用户名或密码错误", c)
		return
	}

	//校验用户
	//dbUser := &models.User{
	//	Username: user.Username,
	//	Password: user.Password,
	//}

	//生成jwt的token
	models.TokenNext(dbUser, c)

}

func getUserInfoAfterLogin(c *gin.Context) {

	userName := c.MustGet(common.GIN_CTX_JWT_USER_NAME).(string)
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	user, err := models.GetUserByUserName(userName)
	if err != nil {
		sc.Logger.Error("通过token解析到的userName去数据库中找User失败",
			zap.Error(err),
		)
		common.ReBadFailWithMessage(fmt.Sprintf("用户名不存在或者密码错误:%v", err.Error()), c)
		return
	}
	common.OkWithDetailed(user, "ok", c)
}
func getPermCode(c *gin.Context) {
	common.OkWithDetailed([]string{"1000", "2000", "3000"}, "ok", c)
}
