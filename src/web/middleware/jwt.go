package middleware

import (
	"bigyunwei-backend/src/common"
	"bigyunwei-backend/src/config"
	"bigyunwei-backend/src/models"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeaderString := c.Request.Header.Get("Authorization")
		if authHeaderString == "" {
			common.Re401FailWithDetailed(gin.H{"reload": true}, "未登录", c)
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeaderString, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			common.Re401FailWithDetailed(gin.H{"reload": true}, "请求头中的auth格式错误", c)
			c.Abort()
			return
		}

		// Parse token
		sc := c.MustGet(common.GIN_CTX_CONFIG_CONFIG).(*config.ServerConfig)
		claims, err := models.ParseToken(parts[1], sc)
		if err != nil {
			common.Re401FailWithDetailed(gin.H{"reload": true}, "解析token出错", c)
			c.Abort()
			return
		}
		// 续期逻辑还没写
		c.Set(common.GIN_CTX_JWT_CLAIM, claims)
		c.Next()
	}
}
