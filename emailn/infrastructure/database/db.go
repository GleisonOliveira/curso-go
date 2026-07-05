package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDb() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		os.Getenv("GORM_POSTGRES_HOST"),
		os.Getenv("GORM_POSTGRES_USER"),
		os.Getenv("GORM_POSTGRES_PASSWORD"),
		os.Getenv("GORM_POSTGRES_DB"),
		os.Getenv("GORM_POSTGRES_PORT"),
		os.Getenv("GORM_POSTGRES_TIMEZONE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Fail to connect to DB")
	}

	return db
}
