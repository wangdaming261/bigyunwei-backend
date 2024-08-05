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
	db, err := gorm.Open(mysql.Open(sc.MysqlC.DSN), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
	})
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
		&Menu{},
	)
}

func MockUserRegister(sc *config.ServerConfig) {
	menus := []*Menu{
		{
			Name:      "System",
			Title:     "系统管理",
			Icon:      "ion:settings-outline",
			Type:      "0",
			Show:      "1",
			OrderNo:   10,
			Component: "LAYOUT",
			Redirect:  "/system/account",
			Path:      "/system",
		},
		{
			Name:       "MenuManagement",
			Title:      "菜单管理",
			Icon:       "ant-design:account-book-filled",
			Type:       "1",
			Show:       "1",
			OrderNo:    11,
			Component:  "/demo/system/menu/index",
			ParentMenu: "1",
			Path:       "menu",
		},
		{
			Name:       "AccountManagement",
			Title:      "用户管理",
			Icon:       "ant-design:account-book-twotone",
			Type:       "1",
			Show:       "1",
			OrderNo:    12,
			Component:  "/demo/system/account/index",
			ParentMenu: "1",
			Path:       "account",
		},
		{
			Name:       "RoleManagement",
			Title:      "角色管理",
			Icon:       "ion:layers-outline",
			Type:       "1",
			Show:       "1",
			OrderNo:    13,
			Component:  "/demo/system/role/index",
			ParentMenu: "1",
			Path:       "role",
		},
		{
			Name:       "ChangePassword",
			Title:      "修改密码",
			Icon:       "ion:layers-outline",
			Type:       "1",
			Show:       "1",
			OrderNo:    14,
			Component:  "/demo/system/password/index",
			ParentMenu: "1",
			Path:       "changePassword",
		},
		{
			Name:      "Permission",
			Title:     "权限管理",
			Icon:      "ion:layers-outline",
			Type:      "0",
			Show:      "1",
			OrderNo:   14,
			Component: "LAYOUT",
			Path:      "/permission",
			Redirect:  "/permission/front/page",
		},
		{
			Name:       "PermissionFrontDemo",
			Title:      "前端权限管理",
			Icon:       "ion:layers-outline",
			Type:       "1",
			Show:       "1",
			OrderNo:    15,
			Component:  "/demo/permission/front/index",
			ParentMenu: "6",
			Path:       "front",
		},
	}

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
				Menus:     menus,
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
