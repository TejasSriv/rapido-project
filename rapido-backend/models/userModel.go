package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Username     string    `gorm:"type:varchar(255);unique;not null" json:"username"`
	Email        string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	FullName     string    `gorm:"type:varchar(255)" json:"fullName"`
	PhoneNumber  string    `gorm:"type:varchar(20)" json:"phoneNumber"`
	Role         string    `gorm:"type:varchar(50);default:'user';not null" json:"role"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}
