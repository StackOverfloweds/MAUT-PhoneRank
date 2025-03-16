package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Battery stores battery-related information
type Battery struct {
	ID                    string `gorm:"type:uuid;primaryKey" json:"id"`
	Capacity              int    `json:"capacity"` // in mAh
	FastChargingAvailable bool   `json:"fast_charging_available"`
	FastCharging          int    `json:"fast_charging"` // in Watts
}

// Generate UUID
func (b *Battery) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String()
	return
}
