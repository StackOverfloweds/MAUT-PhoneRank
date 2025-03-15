package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Smartphone represents main smartphone data, linked to other tables
type Smartphone struct {
	ID                      string    `gorm:"type:uuid;primaryKey" json:"id"`
	BrandID                 string    `json:"brand_id"`
	Brand                   Brand     `gorm:"foreignKey:BrandID" json:"brand"`
	Model                   string    `gorm:"not null" json:"model"`
	Price                   float64   `json:"price"`
	AvgRating               float64   `json:"avg_rating"`
	Is5G                    bool      `json:"is_5G"`
	ProcessorID             string    `json:"processor_id"`
	Processor               Processor `gorm:"foreignKey:ProcessorID" json:"processor"`
	BatteryID               string    `json:"battery_id"`
	Battery                 Battery   `gorm:"foreignKey:BatteryID" json:"battery"`
	DisplayID               string    `json:"display_id"`
	Display                 Display   `gorm:"foreignKey:DisplayID" json:"display"`
	CameraID                string    `json:"camera_id"`
	Camera                  Camera    `gorm:"foreignKey:CameraID" json:"camera"`
	RAMCapacity             int       `json:"ram_capacity"`    // in GB
	InternalMemory          int       `json:"internal_memory"` // in GB
	OS                      string    `json:"os"`
	ExtendedMemoryAvailable bool      `json:"extended_memory_available"`
}

// Generate UUID
func (s *Smartphone) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New().String()
	return
}
