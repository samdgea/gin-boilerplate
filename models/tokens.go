package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type TokenModel struct {
	ID        uuid.UUID `gorm:"primary_key" json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	TokenID   uuid.UUID `json:"token_id"`
	ExpiresAt time.Time `json:"expires_at"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
}

func (base *TokenModel) BeforeCreate(*gorm.DB) (err error) {
	u, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	base.ID = u
	return
}
