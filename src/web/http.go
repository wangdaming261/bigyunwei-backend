package web

import (
	"bigyunwei-backend/src/config"
	"bigyunwei-backend/src/web/view"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func StartGin(sc *config.ServerConfig) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	gin.DisableConsoleColor()
	view.ConfigRoutes(r)
	s := &http.Server{
		Addr:           sc.HttpAddr,
		Handler:        r,
		ReadTimeout:    time.Second * 5,
		WriteTimeout:   time.Second * 5,
		MaxHeaderBytes: 1 << 20,
	}

	//return r.Run(sc.HttpAddr)
	return s.ListenAndServe()
}
