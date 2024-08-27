package models

import (
	"gorm.io/gorm"
)

type UserModelRes struct {
	Id        int    `gorm:"primaryKey" json:"id"`
	Username  string `gorm:"uniqueIndex" json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsActive  bool   `json:"is_active"`
}

type UserModel struct {
	UserModelRes
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func MigrateUserModel(db *gorm.DB) error {
	return db.AutoMigrate(&UserModel{})
}
