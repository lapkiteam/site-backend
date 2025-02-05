package user

import (
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Login    string `gorm:"type:varchar(255);unique"`
	Password string `gorm:"type:varchar(255);not null"`
}
