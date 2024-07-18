package view

import (
	"github.com/gin-gonic/gin"
)

func ConfigRoutes(r *gin.Engine) {
	base := r.Group("/basic-api")
	{
		base.POST("/login", UserLogin)
	}
}
