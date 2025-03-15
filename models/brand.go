package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Brand represents smartphone brands (e.g., Apple, Samsung, etc.)
type Brand struct {
	ID   string `gorm:"type:uuid;primaryKey" json:"id"`
	Name string `gorm:"unique;not null" json:"name"`
}

// Generate UUID before creating a record
func (b *Brand) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String()
	return
}
