package models

import (
	"time"

	"github.com/google/uuid"
)

type Ride struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null" json:"userId"`
	DriverID        *uuid.UUID `gorm:"type:uuid" json:"driverId,omitempty"`
	PickupLocation  string     `gorm:"type:text;not null" json:"pickupLocation"`
	DropoffLocation string     `gorm:"type:text;not null" json:"dropoffLocation"`
	CurrentStatus   string     `gorm:"type:varchar(50);default:'pending';not null" json:"currentStatus"`
	Fare            float64    `gorm:"type:decimal(10,2)" json:"fare"`
	RequestedAt     time.Time  `gorm:"autoCreateTime" json:"requestedAt"`
	AcceptedAt      *time.Time `json:"acceptedAt,omitempty"`
	StartedAt       *time.Time `json:"startedAt,omitempty"`
	CompletedAt     *time.Time `json:"completedAt,omitempty"`
	CancelledAt     *time.Time `json:"cancelledAt,omitempty"`
	AdminNotes      string     `gorm:"type:text" json:"adminNotes,omitempty"`
}

func (Ride) TableName() string {
	return "rides"
}
