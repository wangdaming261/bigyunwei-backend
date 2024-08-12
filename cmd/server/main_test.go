package main

import (
	"bigyunwei-backend/src/common"
	"bigyunwei-backend/src/config"
	"bigyunwei-backend/src/web/middleware"
	"bytes"
	"flag"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"net/http"
	"testing"
	"time"
)

type Menu struct {
	ID         int     `json:"ID"`
	CreatedAt  string  `json:"CreatedAt"`
	UpdatedAt  string  `json:"UpdatedAt"`
	DeletedAt  *string `json:"DeletedAt"`
	Name       string  `json:"name"`
	Title      string  `json:"title"`
	ParentMenu string  `json:"parentMenu"`
	Meta       Meta    `json:"meta"`
	Icon       string  `json:"icon"`
	DbId       int     `json:"DbId"`
	Id         string  `json:"id"`
	Type       string  `json:"type"`
	Show       string  `json:"show"`
	OrderNo    int     `json:"orderNo"`
	Component  string  `json:"component"`
	Redirect   string  `json:"redirect"`
	Path       string  `json:"path"`
	Remark     string  `json:"remark"`
	HomePath   string  `json:"homePath"`
	Status     int     `json:"status"`
	Roles      *string `json:"roles"`
	Children   []Menu  `json:"children"`
}

type Meta struct {
	Title    string `json:"title"`
	Icon     string `json:"icon"`
	ShowMenu bool   `json:"showMenu"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  []Menu `json:"result"`
	Type    string `json:"type"`
}

func TestAdd(t *testing.T) {
	//r := gin.New()
	//r.Use(gin.Recovery())
	var (
		configFile string
	)
	flag.StringVar(&configFile, "config", "D:\\bigyunwei-backend\\server.yml", "path to config file")
	flag.Parse()
	sc, _ := config.LoadServer(configFile)

	logger := common.NewZapLogger(sc.LogLevel, sc.LogFilePath)
	sc.Logger = logger

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
	//r.Use(middleware.NewGinZapLogger(sc.Logger))
	//r.Use(ginzap.Ginzap(sc.Logger, time.RFC3339, false))

	gin.DisableConsoleColor()

	r.Use(ginzap.GinzapWithConfig(sc.Logger, &ginzap.Config{
		Context: func(c *gin.Context) []zapcore.Field {
			var fields []zapcore.Field
			// log request ID
			//if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
			//	fields = append(fields, zap.String("request_id", requestID))
			//}

			if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
				fields = append(fields, zap.String("request_id", requestID))
			}

			// log request body
			var body []byte
			var buf bytes.Buffer
			tee := io.TeeReader(c.Request.Body, &buf)
			body, _ = io.ReadAll(tee)
			c.Request.Body = io.NopCloser(&buf)
			authHeader := c.Request.Header.Get("Authorization")
			fields = append(fields, zap.String("body", string(body)))
			fields = append(fields, zap.String("Authorization", authHeader))

			return fields
		},
	}))

	// Example ping request.
	//r.GET("/ping", func(c *gin.Context) {
	//	c.Writer.Header().Add("X-Request-Id", "1234-5678-9012")
	//	c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	//})
	//
	//r.POST("/ping", func(c *gin.Context) {
	//	c.Writer.Header().Add("X-Request-Id", "9012-5678-1234")
	//	c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	//})
	gin.DisableConsoleColor()
	r.GET("/basic-api/api/system/getMenuList", func(c *gin.Context) {
		c.Writer.Header().Add("X-Request-Id", "9012-5678-1234")

		data := Response{
			Code:    0,
			Message: "ok",
			Result: []Menu{
				{
					ID:         1,
					CreatedAt:  "2024-08-05T17:07:36.426+08:00",
					UpdatedAt:  "2024-08-05T17:07:36.426+08:00",
					DeletedAt:  nil,
					Name:       "System",
					Title:      "系统管理",
					ParentMenu: "",
					Meta: Meta{
						Title:    "系统管理",
						Icon:     "ion:settings-outline",
						ShowMenu: true,
					},
					Icon:      "ion:settings-outline",
					DbId:      0,
					Id:        "1",
					Type:      "0",
					Show:      "1",
					OrderNo:   10,
					Component: "LAYOUT",
					Redirect:  "/system/account",
					Path:      "/system",
					Remark:    "",
					HomePath:  "",
					Status:    0,
					Children: []Menu{
						{
							ID:         2,
							CreatedAt:  "2024-08-05T17:07:36.426+08:00",
							UpdatedAt:  "2024-08-05T17:07:36.426+08:00",
							DeletedAt:  nil,
							Name:       "MenuManagement",
							Title:      "菜单管理",
							ParentMenu: "1",
							Meta: Meta{
								Title:    "菜单管理",
								Icon:     "ant-design:account-book-filled",
								ShowMenu: true,
							},
							Icon:      "ant-design:account-book-filled",
							DbId:      0,
							Id:        "12",
							Type:      "1",
							Show:      "1",
							OrderNo:   11,
							Component: "/demo/system/menu/index",
							Redirect:  "",
							Path:      "menu",
							Remark:    "",
							HomePath:  "",
							Status:    0,
						},
					},
				},
			},
			Type: "",
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    data.Code,
			"message": data.Message,
			"result":  data.Result,
			"type":    data.Type,
		})
	})
	//common.OkWithDetailed(nil, jsonData, c)

	// Listen and Server in 0.0.0.0:8080
	s := &http.Server{
		Addr:           sc.HttpAddr,
		Handler:        r,
		ReadTimeout:    time.Second * 5,
		WriteTimeout:   time.Second * 5,
		MaxHeaderBytes: 1 << 20,
	}

	//return r.Run(sc.HttpAddr)
	_ = s.ListenAndServe()

}
