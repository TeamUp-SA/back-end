package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var gormDB *gorm.DB

func MustInitGorm() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN is empty")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("gorm open: %v", err)
	}
	// Optional: set connection pool on the underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("db.DB(): %v", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	gormDB = db
}

func Gorm() *gorm.DB {
	if gormDB == nil {
		log.Fatal("gorm not initialized; call db.MustInitGorm() in main")
	}
	return gormDB
}
