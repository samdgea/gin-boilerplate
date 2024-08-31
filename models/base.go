package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (base *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	u, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	base.ID = u
	return
}
