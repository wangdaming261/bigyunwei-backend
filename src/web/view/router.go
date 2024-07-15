package view

import (
	"time"

	"github.com/gin-gonic/gin"
)

func ConfigRoutes(r *gin.Engine) {
	base := r.Group("/api")
	{
		base.GET("/ping", ping)
		base.GET("/now", now)
	}
}

func now(r *gin.Context) {
	r.String(200, time.Now().Format("2006-01-02 15:04:05"))
}
