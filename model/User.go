package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if u.Password != nil {
		u.Password, _ = bcrypt.GenerateFromPassword(u.Password, 14)
	}
	return
}
