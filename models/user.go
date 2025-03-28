package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
User - Represents a user in the system.
This model is used for phone authentication.
It includes a unique username, phone number, role, and creation timestamp.
*/
type User struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Phone     string    `gorm:"unique;not null" json:"phone"`
	Role      string    `gorm:"type:varchar(20);default:'user'" json:"role"` // Default role: "user"
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

/*
BeforeCreate - Automatically generates a UUID before creating a new user.
This ensures each user has a unique identifier.
*/
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}
