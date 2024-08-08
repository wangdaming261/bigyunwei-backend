package cron

import (
	"context"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

// 定义同步的manager

func (cm *Manager) SyncCloudResourceManager(ctx context.Context) error {
	go wait.UntilWithContext(ctx, cm.RunSyncCloudResource, time.Duration(cm.Sc.PublicCloudSync.RunIntervalSeconds)*time.Second)
	<-ctx.Done()
	cm.Sc.Logger.Info("SyncCloudResourceManager收到其他任务退出信号 退出")
	return nil

}

func (cm *Manager) RunSyncCloudResource(ctx context.Context) {
	cm.Sc.Logger.Info("模拟同步公有云资源")
}
