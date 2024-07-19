package models

import (
	"bigyunwei-backend/src/common"
	"bigyunwei-backend/src/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
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

func ParseToken(jwtToken string, sc *config.ServerConfig) (*UserCustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		jwtToken,
		&UserCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(sc.JWTC.SigningKey), nil
		},
	)
	if err != nil {
		sc.Logger.Error("解析jwt token失败",
			zap.Error(err))
		return nil, err
	}
	if claims, ok := token.Claims.(*UserCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err

}
