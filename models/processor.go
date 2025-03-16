package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Processor represents CPU details
type Processor struct {
	ID       string  `gorm:"type:uuid;primaryKey" json:"id"`
	Brand    string  `gorm:"not null" json:"brand"`
	Model    string  `gorm:"not null" json:"model"`
	NumCores int     `json:"num_cores"`
	Speed    float64 `json:"speed"` // in GHz
}

// Generate UUID
func (p *Processor) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New().String()
	return
}
