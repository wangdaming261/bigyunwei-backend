package models

// 用户相关的数据库字段

type User struct {
	Username string `json:"username" gorm:"index;comment:用户名"`
	Password string `json:"-" gorm:"comment:密码"`
}
