package models

import (
	"time"

	"github.com/google/uuid"
)

type AdminAction struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	AdminID       uuid.UUID `gorm:"type:uuid;not null" json:"adminId"`
	RideID        uuid.UUID `gorm:"type:uuid;not null" json:"rideId"`
	ActionType    string    `gorm:"type:varchar(50);not null" json:"actionType"`
	ActionDetails string    `gorm:"type:text" json:"actionDetails,omitempty"`
	ActionAt      time.Time `gorm:"autoCreateTime" json:"actionAt"`
	Admin         User      `gorm:"foreignKey:AdminID"`
	Ride          Ride      `gorm:"foreignKey:RideID"`
}

func (AdminAction) TableName() string {
	return "admin_actions"
}
