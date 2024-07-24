package models

import "gorm.io/gorm"

// 用户相关的数据库字段

type User struct {
	gorm.Model
	UserId   int    `json:"userId" gorm:"comment:用户ID"`
	Username string `json:"username" gorm:"index;comment:用户名"`
	Password string `json:"password" gorm:"comment:密码"`
	RealName string `json:"realName" gorm:"comment:昵称"`
	Desc     string `json:"desc" gorm:"comment:描述"`
	HomePath string `json:"homePath" gorm:"comment:主页"`
	Enable   int    `json:"enable" gorm:"default:1;comment:是否启用 1启用 0禁用"`
	//Roles    []*Role `json:"roles" gorm:"many2many:user_roles;comment:用户角色"`
	//Phone    string  `json:"phone" gorm:"comment:手机号"`
	//Email    string  `json:"email" gorm:"comment:邮箱"`
}
