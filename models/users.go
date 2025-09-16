package models

import "golang.org/x/crypto/bcrypt"

type Users struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Email    string `json:"email" gorm:"unique"`
	AddTime  string `json:"add_time"`
}

// 配置操作数据库的表名称
func (Users) TableName() string {
	return "users"
}

// 密码加密
func (user *Users) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// 验证
func (user *Users) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
