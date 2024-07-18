package models

import (
	"bigyunwei-backend/src/common"
	"bigyunwei-backend/src/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"time"
)

func TokenNext(dbUser *User, c *gin.Context) {
	//生成jwt的token
	sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
	token, err := GenJwtToken(dbUser, sc)
	if err != nil {
		sc.Logger.Error("生成jwt token失败", zap.Error(err))
		common.FailWithMessage("生成jwt token失败", c)
		return
	}
	//返回token
	common.OkWithDetailed(UserLoginResponse{
		User:  dbUser,
		Token: token,
	}, "登录成功", c)

}

func GenJwtToken(dbUser *User, sc *config.ServerConfig) (string, error) {
	c := UserCustomClaims{
		User: dbUser,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(sc.JWTC.ExpireDuration)),
			Issuer:    sc.JWTC.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(sc.JWTC.SigningKey))
}
