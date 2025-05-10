package main

import (
	"fmt"
	"github.com/SwanHtetAungPhyo/supabase-migration/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	//dsn := "host=database-1.cm9ewocwci8f.us-east-1.rds.amazonaws.com user=postgres password=Swanhtetaungphyo dbname=postgres port=5432 sslmode=require"
	dsn := "host=localhost user=postgres password=postgres dbname=auth port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}
	fmt.Println("✅ Connected to RDS")

	// Auto-migrate all tables (drop Chat/Message, but include Order)
	if err := db.AutoMigrate(
		&model.Badge{},
		&model.UserBadge{},
		&model.User{},
		&model.Skill{},
		&model.UserSkill{},
		&model.Biometrics{},
		&model.GigTag{},
		&model.Gig{},
		&model.RegistrationToken{},
		&model.GigImage{},
		&model.GigPackage{},
		&model.GigPackageFeature{},
		&model.Category{},
		&model.Order{}, // <-- ensures Order table is created/updated
		&model.Review{},
		//&model.GigToGigTag{},
	); err != nil {
		fmt.Println("❌ AutoMigrate failed:", err)
		return
	}
	fmt.Println("✅ All models migrated (including Order)")
}
