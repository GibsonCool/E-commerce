package datamodels

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	NickName string
	UserName string `gorm:"unique;not null"` //唯一，切不能为空
	HashPwd  string `gorm:"unique"`
}

