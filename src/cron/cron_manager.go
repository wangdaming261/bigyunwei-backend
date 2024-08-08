package cron

import "bigyunwei-backend/src/config"

type Manager struct {
	Sc *config.ServerConfig
}

func NewManager(sc *config.ServerConfig) *Manager {
	return &Manager{Sc: sc}
}
