package database

import (
	"fmt"
	"github.com/SwanHtetAungPhyo/migration_service/providers/model"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func NewDB() (*gorm.DB, error) {
	dsn := viper.GetString("sup-abase.url")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.Country{},
		&model.User{},
		&model.Order{},
		&model.Chat{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	fmt.Println("✅ Database migration completed.")
	return nil
}
func isAlreadyExistsError(err error) bool {
	return err != nil && (err.Error() == "relation already exists" ||
		err.Error() == "constraint already exists")
}
