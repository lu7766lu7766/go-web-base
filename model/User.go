package model

import (
	"time"
)

type User struct {
	Id       uint   `gorm:"column:id;primaryKey;type:AUTO_INCREMENT" json:"id"`
	Password []byte `json:"-"`
	Name     string `json:"name"`
	Mail     string `json:"mail"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// func (User) TableName() string {
// 	return "users"
// }

// func (u *User) AfterFind(tx *gorm.DB) (err error) {
// 	return
// }
