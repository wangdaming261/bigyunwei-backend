package main

import (
	"bigyunwei-backend/src/common"
	"bigyunwei-backend/src/config"
	"bigyunwei-backend/src/cron"
	"bigyunwei-backend/src/models"
	"bigyunwei-backend/src/web"
	"context"
	"flag"
	"fmt"

	esl "github.com/ning1875/errgroup-signal/signal"
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

	logger := common.NewZapLogger(sc.LogLevel, sc.LogFilePath)
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("关闭日志失败: %v\n", err)
		}
	}(logger)
	sc.Logger = logger
	// 初始化数据库
	err = models.InitDB(sc)
	if err != nil {
		sc.Logger.Error("初始化数据库失败", zap.Error(err))
		return
	}
	err = models.MigrateTable()
	if err != nil {
		sc.Logger.Error("迁移表失败", zap.Error(err))
		return
	}
	// TODO 测试用，后续可以删除
	models.MockUserRegister(sc)

	cm := cron.NewManager(sc)

	group, stopChan := esl.SetupStopSignalContext()
	ctxAll, cancelAll := context.WithCancel(context.Background())

	group.Go(func() error {
		logger.Info("[stop chan watch start backend]")
		for {
			select {
			case <-stopChan:
				logger.Info("stop chan receive quite signal exit")
				cancelAll()
				return nil
			}
		}
	})

	group.Go(func() error {
		logger.Info("计划任务--同步公有云启动")
		err := cm.SyncCloudResourceManager(ctxAll)
		if err != nil {
			logger.Error("计划任务--同步公有云失败", zap.Error(err))
		}
		return err
	})

	group.Go(func() error {

		errChan := make(chan error, 1)
		go func() {
			errChan <- web.StartGin(sc)
		}()
		logger.Info("start backend")
		select {
		case err := <-errChan:
			logger.Error("[]", zap.Error(err))
			return err
		case <-ctxAll.Done():
			return nil
		}
	})
	group.Wait()
	//err = web.StartGin(sc)
	//if err != nil {
	//	return
	//}
}
