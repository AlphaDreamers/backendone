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
	// Disable foreign key checks temporarily
	if err := db.Exec("SET session_replication_role = 'replica'").Error; err != nil {
		return fmt.Errorf("failed to disable foreign key checks: %w", err)
	}

	// Create tables without foreign keys first
	tables := []struct {
		name  string
		model interface{}
	}{
		{"users", &model.User{}},
		{"service_posts", &model.ServicePost{}},
		{"firebase_tokens", &model.FirebaseToken{}},
		{"reviews", &model.Review{}},
	}

	for _, table := range tables {
		log.Printf("Creating table: %s", table.name)
		if err := db.Migrator().CreateTable(table.model); err != nil {
			if !isAlreadyExistsError(err) {
				return fmt.Errorf("failed to create table %s: %w", table.name, err)
			}
			log.Printf("Table %s already exists, skipping creation", table.name)
		}
	}

	// Add constraints only if they don't exist
	constraints := []struct {
		name       string
		sql        string
		checkQuery string
	}{
		{
			"fk_users_firebase_token",
			`ALTER TABLE firebase_tokens 
			 ADD CONSTRAINT fk_users_firebase_token 
			 FOREIGN KEY (user_id) REFERENCES users(id) 
			 ON DELETE CASCADE`,
			`SELECT 1 FROM pg_constraint WHERE conname = 'fk_users_firebase_token'`,
		},
		{
			"fk_service_posts_owner",
			`ALTER TABLE service_posts 
			 ADD CONSTRAINT fk_service_posts_owner 
			 FOREIGN KEY (owner_id) REFERENCES users(id) 
			 ON DELETE CASCADE`,
			`SELECT 1 FROM pg_constraint WHERE conname = 'fk_service_posts_owner'`,
		},
		{
			"fk_reviews_reviewer",
			`ALTER TABLE reviews 
			 ADD CONSTRAINT fk_reviews_reviewer 
			 FOREIGN KEY (reviewer_id) REFERENCES users(id) 
			 ON DELETE SET NULL`,
			`SELECT 1 FROM pg_constraint WHERE conname = 'fk_reviews_reviewer'`,
		},
		{
			"fk_service_posts_reviews",
			`ALTER TABLE reviews 
			 ADD CONSTRAINT fk_service_posts_reviews 
			 FOREIGN KEY (service_id) REFERENCES service_posts(service_id) 
			 ON DELETE CASCADE`,
			`SELECT 1 FROM pg_constraint WHERE conname = 'fk_service_posts_reviews'`,
		},
	}

	for _, constraint := range constraints {
		var exists int
		if err := db.Raw(constraint.checkQuery).Scan(&exists).Error; err != nil {
			return fmt.Errorf("failed to check constraint %s: %w", constraint.name, err)
		}

		if exists == 0 {
			log.Printf("Adding constraint: %s", constraint.name)
			if err := db.Exec(constraint.sql).Error; err != nil {
				return fmt.Errorf("failed to add constraint %s: %w", constraint.name, err)
			}
		} else {
			log.Printf("Constraint %s already exists, skipping", constraint.name)
		}
	}

	// Re-enable foreign key checks
	if err := db.Exec("SET session_replication_role = 'origin'").Error; err != nil {
		return fmt.Errorf("failed to enable foreign key checks: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

func isAlreadyExistsError(err error) bool {
	return err != nil && (err.Error() == "relation already exists" ||
		err.Error() == "constraint already exists")
}
