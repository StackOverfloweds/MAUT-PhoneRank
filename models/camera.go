package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Camera holds smartphone camera details
type Camera struct {
	ID                 string  `gorm:"type:uuid;primaryKey" json:"id"`
	NumRearCameras     int     `json:"num_rear_cameras"`
	PrimaryCameraRear  float64 `json:"primary_camera_rear"`  // in MP
	PrimaryCameraFront float64 `json:"primary_camera_front"` // in MP
}

// Generate UUID
func (c *Camera) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()
	return
}
