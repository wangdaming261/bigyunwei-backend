package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	OrderNo   int     `json:"orderNo" gorm:"comment:'排序'"`
	RoleName  string  `json:"roleName" gorm:"type:varchar(100);uniqueIndex;comment:'角色中文名称'"` // 用户登录名 index 代表索引
	RoleValue string  `json:"roleValue" gorm:"type:varchar(100);uniqueIndex;comment:'角色值'"`
	Remark    string  `json:"remark" gorm:"comment:'用户描述'"`
	HomePath  string  `json:"HomePath" gorm:"comment:'登录后的默认访问页'"`
	Status    int     `json:"status" gorm:"default:1;comment:'角色是否被冻结 1正常 2冻结'"` // 用户是否被冻结 1正常 2冻结
	Users     []*User `gorm:"many2many:user_roles;comment:'用户'"`
}
