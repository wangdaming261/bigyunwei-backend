package models

import (
	"bigyunwei-backend/src/common"
	"bigyunwei-backend/src/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Db *gorm.DB
)

func InitDB(sc *config.ServerConfig) error {
	db, err := gorm.Open(mysql.Open(sc.MysqlC.DSN), &gorm.Config{})
	if err != nil {
		return err
	}
	Db = db
	return nil

}

func MigrateTable() error {
	return Db.AutoMigrate(
		&User{},
		&Role{},
	)
}

func MockUserRegister(sc *config.ServerConfig) {
	u1 := User{
		Username: "admin",
		Password: "123456",
		RealName: "管理员",
		Desc:     "",
		HomePath: "/system/account",
		Enable:   1,
		Roles: []*Role{
			{
				RoleName:  "超级管理员",
				RoleValue: "super",
			},
			{
				RoleName:  "前端管理员",
				RoleValue: "frontAdmin",
			},
		},
	}

	u1.Password = common.BcryptHash(u1.Password)
	err := Db.Create(&u1).Error
	if err != nil {
		sc.Logger.Error("模拟用户注册失败", zap.String("错误", err.Error()))
		return
	}

	sc.Logger.Info("模拟注册成功")
}
