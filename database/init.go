package database

import (
	"fmt"
	"log"
	"os"

	migrate "github.com/StackOverfloweds/MAUT-PhoneRank/database/Migrate"
	csv "github.com/StackOverfloweds/MAUT-PhoneRank/database/csv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	fmt.Println("Database connected successfully!")

	migrate.MigrateDB(db)

	csv.ImportSmartphones(db)
}
