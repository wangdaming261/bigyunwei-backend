package cron

import (
	"bigyunwei-backend/src/config"
	"sync"
)

type Manager struct {
	Sc *config.ServerConfig
	//=true 代表上次已经同步完了。这次可以:担心上次没执行定，新的就开始了:数据冲突:最多只有1个同步任务
	EcsLastSyncFinished bool
	RdsLastSyncFinished bool
	ElbLastSyncFinished bool
	sync.RWMutex
}

func (cm *Manager) SetEcsSynced(fin bool) {
	cm.Lock()
	defer cm.Unlock()
	cm.EcsLastSyncFinished = fin
}

func (cm *Manager) GetEcsSynced() bool {
	cm.RLock()
	defer cm.RUnlock()
	return cm.EcsLastSyncFinished
}

func NewManager(sc *config.ServerConfig) *Manager {
	return &Manager{Sc: sc,
		EcsLastSyncFinished: true,
	}
}
