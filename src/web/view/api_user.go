package view

import (
	"bigyunwei-backend/src/common"
	"bigyunwei-backend/src/config"
	"bigyunwei-backend/src/models"
	"errors"
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

func getUserInfoAfterLog(c *gin.Context) {
	jwtClaim := c.MustGet(common.GIN_CTX_JWT_CLAIM).(*models.UserCustomClaims)
	common.OkWithDetailed(jwtClaim.User, "ok", c)
}
