package models

import (
	"bigyunwei-backend/src/config"
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
	)
}
