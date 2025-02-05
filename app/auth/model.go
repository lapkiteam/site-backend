package auth

import "gorm.io/gorm"

type SessionModel struct {
	gorm.Model
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Token string `gorm:"type:text"`
}
