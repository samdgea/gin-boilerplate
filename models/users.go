package models

import (
	"gorm.io/gorm"
)

type UserModel struct {
	BaseModel
	Username  string `gorm:"uniqueIndex" json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsActive  bool   `json:"is_active" gorm:"default:false"`
	Password  string `json:"password"`
}

func MigrateUserModel(db *gorm.DB) error {
	return db.AutoMigrate(&UserModel{})
}
