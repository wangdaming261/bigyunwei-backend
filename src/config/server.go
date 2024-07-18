package config

import (
	"os"
	"time"

	"go.uber.org/zap"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	HttpAddr    string      `yaml:"http_addr"`
	LogLevel    string      `yaml:"log_level"`
	LogFilePath string      `yaml:"log_file_path"`
	JWTC        *JWT        `yaml:"jwt"`
	Logger      *zap.Logger `yaml:"-"`
}

func LoadServer(filename string) (*ServerConfig, error) {
	cfg := &ServerConfig{}
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		return nil, err
	}
	// 解析过期时间
	duration, err := time.ParseDuration(cfg.JWTC.ExpireTime)
	if err != nil {
		return nil, err
	}
	cfg.JWTC.ExpireDuration = duration
	return cfg, err
}

type JWT struct {
	SigningKey     string        `yaml:"signing_key" json:"signing_key"`
	ExpireTime     string        `yaml:"expire_time" json:"expire-time"`
	ExpireDuration time.Duration `yaml:"-"`
	BufferTime     string        `yaml:"buffer_time" json:"buffer-time"`
	Issuer         string        `yaml:"issuer" json:"issuer"`
}
