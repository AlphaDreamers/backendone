package db

import (
	"log"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

func InitDB(dsn string) {
	once.Do(func() {
		var err error
		var tempDB *gorm.DB

		// Retry logic: Try connecting up to 5 times
		for i := 0; i < 5; i++ {
			tempDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info), // Enable logging
			})
			if err == nil {
				break
			}
			log.Printf("Database connection failed (attempt %d/5): %v", i+1, err)
			time.Sleep(2 * time.Second)
		}

		db = tempDB
		sqlDB, err := db.DB()
		if err != nil {
			log.Println("Cannot set the connection pool")
		}
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxIdleTime(25 * time.Minute)
		log.Println("Database connected successfully!")
	})
}

func GetDB() *gorm.DB {
	return db
}
