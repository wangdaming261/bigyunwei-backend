package web

import (
	"bigyunwei-backend/src/common"
	"bigyunwei-backend/src/config"
	"bigyunwei-backend/src/web/middleware"
	"bigyunwei-backend/src/web/view"
	"net/http"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func StartGin(sc *config.ServerConfig) error {
	//gin.SetMode(gin.ReleaseMode)
	//r := gin.Default()
	r := gin.New()
	r.Use(gin.Recovery())

	varMap := map[string]interface{}{}
	//varMap[common.GIN_CTX_CONFIG_LOGGER] = sc.Logger
	varMap[common.GIN_CTX_CONFIG_CONFIG] = sc
	// 传递变量
	r.Use(middleware.ConfigMiddleware(varMap))
	// 打印耗时
	//r.Use(middleware.TimeCost())
	// 请求id
	r.Use(requestid.New())
	// 自定义日志
	r.Use(middleware.NewGinZapLogger(sc.Logger))
	//r.Use(ginzap.Ginzap(sc.Logger, time.RFC3339, false))

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
