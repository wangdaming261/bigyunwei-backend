package main

import (
	"bigyunwei-backend/src/common"
	"bigyunwei-backend/src/config"
	"bigyunwei-backend/src/web"
	"flag"
	"fmt"

	"go.uber.org/zap"
)

func main() {
	var (
		configFile string
	)

	flag.StringVar(&configFile, "config", "./server.yml", "path to config file")
	flag.Parse()

	sc, err := config.LoadServer(configFile)
	if err != nil {
		fmt.Printf("加载配置文件失败: %v\n", err)
		panic(err)
	}
	fmt.Printf("主配置文件路径: %v sc: %v\n", configFile, sc)

	logger := common.NewZapLogger(sc.LogLevel)
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("关闭日志失败: %v\n", err)
		}
	}(logger)
	logger.Debug("服务启动")

	err = web.StartGin(sc)
	if err != nil {
		return
	}

}
