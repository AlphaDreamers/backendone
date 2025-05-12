package main

import (
	"fmt"
	"github.com/SwanHtetAungPhyo/supabase-migration/model"
	"github.com/google/uuid"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func Connect() {
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

//func main() {
//	dsn := "host=localhost user=postgres password=postgres dbname=auth port=5432 sslmode=disable"
//	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//	if err != nil {
//		panic("failed to connect to database: " + err.Error())
//	}
//	fmt.Println("✅ Connected to RDS")
//}

func InitializeDatabase(db *gorm.DB) error {
	// Start a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Create independent entities first
	badges := []model.Badge{
		{Label: "Top Seller", Icon: "medal", Color: "gold"},
		{Label: "Fast Delivery", Icon: "rocket", Color: "blue"},
		{Label: "5-Star Rating", Icon: "star", Color: "purple"},
	}
	if err := tx.Create(&badges).Error; err != nil {
		tx.Rollback()
		return err
	}

	skills := []model.Skill{
		{Label: "Web Development"},
		{Label: "Graphic Design"},
		{Label: "Digital Marketing"},
	}
	if err := tx.Create(&skills).Error; err != nil {
		tx.Rollback()
		return err
	}

	gigTags := []model.GigTag{
		{Label: "Frontend"},
		{Label: "Logo Design"},
		{Label: "SEO"},
	}
	if err := tx.Create(&gigTags).Error; err != nil {
		tx.Rollback()
		return err
	}

	categories := []model.Category{
		{Label: "Digital Services", Slug: "digital-services", IsActive: true, SortOrder: 1},
		{Label: "Creative Design", Slug: "creative-design", IsActive: true, SortOrder: 2},
		{Label: "Marketing", Slug: "marketing", IsActive: true, SortOrder: 3},
	}
	if err := tx.Create(&categories).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. Create users
	users := []model.User{
		{
			FirstName:       "John",
			LastName:        "Doe",
			Email:           "john.doe@example.com",
			CognitoUsername: uuid.New().String(),
			Username:        "johndoe",
			Country:         "US",
		},
		{
			FirstName:       "Jane",
			LastName:        "Smith",
			Email:           "jane.smith@example.com",
			CognitoUsername: uuid.New().String(),
			Username:        "janesmith",
			Country:         "UK",
		},
	}
	if err := tx.Create(&users).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 3. Add user skills and badges
	userSkills := []model.UserSkill{
		{
			UserID:  users[0].ID,
			SkillID: skills[0].ID, // Web Development
			Level:   3,
		},
		{
			UserID:  users[1].ID,
			SkillID: skills[1].ID, // Graphic Design
			Level:   2,
		},
	}
	if err := tx.Create(&userSkills).Error; err != nil {
		tx.Rollback()
		return err
	}

	userBadges := []model.UserBadge{
		{
			UserID:     users[0].ID,
			BadgeID:    badges[0].ID, // Top Seller
			Tier:       "GOLD",
			IsFeatured: true,
		},
		{
			UserID:     users[1].ID,
			BadgeID:    badges[1].ID, // Fast Delivery
			Tier:       "SILVER",
			IsFeatured: false,
		},
	}
	if err := tx.Create(&userBadges).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 4. Create gig with first user as seller
	gig := model.Gig{
		Title:       "Professional Website Development",
		Description: "I will build you a responsive website",
		CategoryID:  categories[0].ID, // Digital Services
		SellerID:    users[0].ID,
		Tags: []model.GigTag{
			gigTags[0], // Frontend
		},
		Packages: []model.GigPackage{
			{
				Title:        "Basic Website",
				Description:  "5-page responsive website",
				Price:        500,
				DeliveryTime: 14,
				Features: []model.GigPackageFeature{
					{Title: "Responsive Design", Included: true},
					{Title: "Contact Form", Included: true},
				},
			},
			{
				Title:        "Premium Website",
				Description:  "10-page website with CMS",
				Price:        1000,
				DeliveryTime: 21,
				Features: []model.GigPackageFeature{
					{Title: "Content Management System", Included: true},
					{Title: "SEO Optimization", Included: true},
				},
			},
		},
		Images: []model.GigImage{
			{URL: "https://example.com/gig1.jpg", IsPrimary: true},
			{URL: "https://example.com/gig2.jpg"},
		},
	}
	if err := tx.Create(&gig).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 5. Create an order with second user as buyer
	now := time.Now()
	order := model.Order{
		OrderNumber:   "ORD-001",
		Price:         gig.Packages[0].Price,
		PaymentMethod: "credit_card",
		PackageID:     gig.Packages[0].ID,
		SellerID:      users[0].ID,
		BuyerID:       users[1].ID,
		Status:        "COMPLETED",
		CompletedAt:   &now,
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return err
	}

	reveiwContent := "The website looks amazing and was delivered on time"
	// 6. Add review for the order
	review := model.Review{
		Title:    "Excellent work!",
		Content:  &reveiwContent,
		Rating:   5,
		AuthorID: users[1].ID,
		OrderID:  order.ID,
	}
	if err := tx.Create(&review).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

func main() {
	// Initialize your database connection
	dsn := "host=localhost user=postgres password=postgres dbname=auth port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run the initialization
	if err := InitializeDatabase(db); err != nil {
		log.Fatal("Database initialization failed:", err)
	}

	log.Println("Database initialized successfully!")
}
