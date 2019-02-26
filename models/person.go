package models


import (
	"github.com/jinzhu/gorm"
)


// 参考 https://www.jianshu.com/p/443766f0e796

// db model
type Person struct {
	//ID uint `json:"id"`
	gorm.Model
	Firstname string `json:"firstname" binding:"required"`
	Lastname  string `json:"lastname" binding:"required"`
	City      string `json:"city" binding:"required"`
}

// 自定义表名
func (Person) TableName() string {
	return "persons"
}


