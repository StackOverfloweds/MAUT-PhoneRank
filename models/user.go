package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User table for authentication
type User struct {
	ID           string    `gorm:"type:uuid;primaryKey" json:"id"`
	Username     string    `gorm:"unique;not null" json:"username"`
	Email        string    `gorm:"unique;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"` // Hide password in JSON response
	Role         string    `gorm:"default:user" json:"role"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	Profile      Profile   `gorm:"foreignKey:UserID" json:"profile"`
}

// Auto-generate UUID before creating a new user
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}
