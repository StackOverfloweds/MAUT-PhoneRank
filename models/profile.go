package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Profile table for user identity
type Profile struct {
	ID          string     `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      string     `gorm:"unique;not null" json:"user_id"`
	FullName    string     `gorm:"not null" json:"full_name"`
	Address     *string    `json:"address,omitempty"`                                                                      // Nullable
	PhoneNumber *string    `gorm:"unique" json:"phone_number,omitempty"`                                                   // Nullable
	Birthdate   *time.Time `json:"birthdate,omitempty"`                                                                    // Nullable
	Gender      *string    `gorm:"type:varchar(10);check:(gender IN ('Male', 'Female', 'Other'))" json:"gender,omitempty"` // Nullable
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// Auto-generate UUID before creating a new profile
func (p *Profile) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New().String()
	return
}
