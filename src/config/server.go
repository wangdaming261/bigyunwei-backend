package config

import (
	"os"
	"time"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
)

type ServerConfig struct {
	HttpAddr    string        `yaml:"http_addr"`
	MysqlC      *mysql.Config `yaml:"mysql"`
	LogLevel    string        `yaml:"log_level"`
	LogFilePath string        `yaml:"log_file_path"`
	//SuperRoleName  string           `yaml:"super_role_name"`
	PublicCloudSync *PublicCloudSync `yaml:"public_cloud_sync"`
	JWTC            *JWT             `yaml:"jwt"`
	Logger          *zap.Logger      `yaml:"-"`
}

type PublicCloudSync struct {
	RunIntervalSeconds int         `yaml:"run_interval_seconds"`
	AliCloud           []*AliCloud `yaml:"ali_cloud"`
}

type AliCloud struct {
	Enable          bool   `yaml:"enable"`
	RegionId        string `yaml:"region_id"`
	AccessKeyId     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
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

	// 解析缓冲时间
	duration, err = time.ParseDuration(cfg.JWTC.BufferTime)
	if err != nil {
		return nil, err

	}
	cfg.JWTC.BufferDuration = duration

	return cfg, err
}

type JWT struct {
	SigningKey     string        `yaml:"signing_key" json:"signing_key"` //jwt签名 密码加盐
	ExpireTime     string        `yaml:"expire_time" json:"expire_time"` //过期时间
	ExpireDuration time.Duration `yaml:"-"`                              //解析配置文件用
	BufferTime     string        `yaml:"buffer_time" json:"buffer_time"` //过期时间
	BufferDuration time.Duration `yaml:"-"`                              //临期时间
	Issuer         string        `yaml:"issuer" json:"issuer"`           //签发者
}
