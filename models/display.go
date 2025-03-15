package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Display represents screen specifications
type Display struct {
	ID               string  `gorm:"type:uuid;primaryKey" json:"id"`
	ScreenSize       float64 `json:"screen_size"`  // in inches
	RefreshRate      int     `json:"refresh_rate"` // in Hz
	ResolutionWidth  int     `json:"resolution_width"`
	ResolutionHeight int     `json:"resolution_height"`
}

// Generate UUID
func (d *Display) BeforeCreate(tx *gorm.DB) (err error) {
	d.ID = uuid.New().String()
	return
}
