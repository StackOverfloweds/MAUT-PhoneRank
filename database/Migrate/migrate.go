package migrate

import (
	"github.com/StackOverfloweds/MAUT-PhoneRank/models"
	"gorm.io/gorm"
)

// MigrateDB creates necessary tables
func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Brand{},
		&models.Processor{},
		&models.Battery{},
		&models.Display{},
		&models.Camera{},
		&models.Smartphone{},
	)
}
