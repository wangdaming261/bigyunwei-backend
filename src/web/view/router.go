package view

import (
	"bigyunwei-backend/src/web/middleware"

	"github.com/gin-gonic/gin"
)

func ConfigRoutes(r *gin.Engine) {
	base := r.Group("/basic-api")
	{
		base.POST("/login", UserLogin)
	}

	afterLoginApiGroup := r.Group("/api")
	afterLoginApiGroup.Use(middleware.JWTAuthMiddleware())
	{
		afterLoginApiGroup.GET("/user/info", getUserInfoAfterLog)
	}
}
