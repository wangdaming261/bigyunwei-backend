package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Menu struct {
	gorm.Model
	Name       string    `json:"name" gorm:"type:varchar(100);uniqueIndex;comment:名称"`
	Title      string    `json:"title" gorm:"comment:名称"`
	ParentMenu string    `json:"parentMenu" gorm:"type:varchar(5);comment:父级的id"`
	Meta       *MenuMeta `json:"meta" gorm:"-"`
	Icon       string    `json:"icon" gorm:"comment:图标"`
	DbId       uint      `json:"DbId" gorm:"-"`
	Id         string    `json:"id" gorm:"-"`
	Type       string    `json:"type" gorm:"type:varchar(5);comment:类型 0=目录 1=子菜单"`
	Show       string    `json:"show" gorm:"type:varchar(5);comment:类型 0=禁用 1=启用"`
	OrderNo    int       `json:"orderNo" gorm:"comment:排序"`
	Component  string    `json:"component" gorm:"type:varchar(50);comment:前端组件 菜单就是LAYOUT"`
	Redirect   string    `json:"redirect" gorm:"type:varchar(50);comment:显示路径"`
	Path       string    `json:"path" gorm:"type:varchar(50);comment:路由路径"`
	Remark     string    `json:"remark" gorm:"comment:用户描述"`
	HomePath   string    `json:"homePath" gorm:"comment:登录后的默认首页"`
	Status     int       `json:"status" gorm:"default:0;comment:角色是否被禁用 0正常 1禁用"`
	Roles      []*Role   `json:"roles" gorm:"many2many:role_menus;"`
	Children   []*Menu   `json:"children" gorm:"-"`
}

type MenuMeta struct {
	Title    string `json:"title" gorm:"-"`    // 标题
	Icon     string `json:"icon" gorm:"-"`     // 图标
	ShowMenu bool   `json:"showMenu" gorm:"-"` // 是否显示
}

func GetMenuById(id int) (*Menu, error) {
	var dbMenu Menu
	err := Db.Where("id = ?", id).Preload("Roles").First(&dbMenu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("菜单不存在")
		}
		return nil, fmt.Errorf("数据库错误: %w", err)
	}
	return &dbMenu, nil
}
