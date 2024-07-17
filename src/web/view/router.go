package view

import (
	"time"

	"github.com/gin-gonic/gin"
)

func ConfigRoutes(r *gin.Engine) {
	base := r.Group("/basic-api")
	{
		base.GET("/ping", ping)
		base.GET("/now", now)
		base.POST("/login", UserLogin)
	}
}

func now(r *gin.Context) {
	//time.Sleep(3 * time.Second)
	r.String(200, time.Now().Format("2006-01-02 15:04:05"))
}
