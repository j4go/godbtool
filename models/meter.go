package models

import (
	"github.com/jinzhu/gorm"
)

// 参考 https://www.jianshu.com/p/443766f0e796

// db model
type Meter struct {
	//ID uint `json:"id"`
	gorm.Model
	MeterId string `json:"meterid" binding:"required" gorm:"unique_index:hash_idx;"`
	Config  string `json:"config" binding:"required" `
}

// 自定义表名
func (Meter) TableName() string {
	return "meter_config"
}
