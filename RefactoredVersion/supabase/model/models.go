package model

import (
	"github.com/google/uuid"
	"time"
)

// Badge maps to the "Badge" table
// PK: id TEXT
// Columns: label, icon, color, createdAt, updatedAt
// Relations: none

type Badge struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Label     string    `gorm:"column:label;type:text;not null"`
	Icon      string    `gorm:"column:icon;type:text;not null"`
	Color     string    `gorm:"column:color;type:text;not null"`
	CreatedAt time.Time `gorm:"column:createdAt;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updatedAt;autoUpdateTime"`
}

func (Badge) TableName() string { return "Badge" }

// UserBadge maps to the "UserBadge" table
// PK: id TEXT
// FKs: userId -> User.id, badgeId -> Badge.id

type UserBadge struct {
	ID         uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID `gorm:"column:userId;type:uuid;not null"`
	BadgeID    uuid.UUID `gorm:"column:badgeId;type:uuid;not null"`
	Tier       string    `gorm:"column:tier;type:text;not null;default:'BRONZE'"`
	IsFeatured bool      `gorm:"column:isFeatured;not null;default:false"`
	CreatedAt  time.Time `gorm:"column:createdAt;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"column:updatedAt;autoUpdateTime"`

	User  User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Badge Badge `gorm:"foreignKey:BadgeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (UserBadge) TableName() string { return "UserBadge" }

// User maps to the "User" table

type User struct {
	ID                uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	FirstName         string     `gorm:"column:firstName;type:text;not null"`
	LastName          string     `gorm:"column:lastName;type:text;not null"`
	Email             string     `gorm:"column:email;type:text;not null;uniqueIndex"`
	CognitoUsername   string     `gorm:"column:cognito_user_name;type:uuid;not null"`
	Verified          bool       `gorm:"column:verified;not null;default:false"`
	TwoFactorVerified bool       `gorm:"column:twoFactorVerified;not null;default:false"`
	Username          string     `gorm:"column:username;type:text;not null;"`
	Avatar            *string    `gorm:"column:avatar;type:text"`
	Country           string     `gorm:"column:country;type:text;not null"`
	WalletCreated     bool       `gorm:"column:walletCreated;not null;default:false"`
	WalletCreatedTime *time.Time `gorm:"column:walletCreatedTime"`
	CreatedAt         time.Time  `gorm:"column:createdAt;autoCreateTime"`
	UpdatedAt         time.Time  `gorm:"column:updatedAt;autoUpdateTime"`

	UserBadges     []UserBadge  `gorm:"foreignKey:UserID"`
	Skills         []UserSkill  `gorm:"foreignKey:UserID"`
	Biometrics     []Biometrics `gorm:"foreignKey:UserID"`
	Gigs           []Gig        `gorm:"foreignKey:SellerID"`
	OrdersAsBuyer  []Order      `gorm:"foreignKey:BuyerID"`
	OrdersAsSeller []Order      `gorm:"foreignKey:SellerID"`
}

func (User) TableName() string { return "User" }

// Skill maps to the "Skill" table

type Skill struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Label     string    `gorm:"column:label;type:text;not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"column:createdAt;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updatedAt;autoUpdateTime"`

	Users []UserSkill `gorm:"foreignKey:SkillID"`
}

func (Skill) TableName() string { return "Skill" }

// UserSkill maps to the "user_skills" join table

type UserSkill struct {
	SkillID   uuid.UUID `gorm:"column:skillId;type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"column:userId;type:uuid;primaryKey"`
	Level     int       `gorm:"column:level;not null;default:1"`
	Endorsed  bool      `gorm:"column:endorsed;not null;default:false"`
	CreatedAt time.Time `gorm:"column:createdAt;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updatedAt;autoUpdateTime"`

	User  User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Skill Skill `gorm:"foreignKey:SkillID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (UserSkill) TableName() string { return "user_skills" }

// Biometrics maps to the "Biometrics" table

type Biometrics struct {
	ID              uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	CognitoUsername *string   `gorm:"column:cognito_user_name;type:uuid;not null"`
	Type            string    `gorm:"column:type;type:text;not null"`
	Value           string    `gorm:"column:value;type:text;not null"`
	IsVerified      bool      `gorm:"column:isVerified;not null;default:false"`
	UserID          uuid.UUID `gorm:"column:userId;type:uuid;not null"`
	CreatedAt       time.Time `gorm:"column:createdAt;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"column:updatedAt;autoUpdateTime"`

	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Biometrics) TableName() string { return "Biometrics" }

// GigTag maps to the "GigTag" table

type GigTag struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Label     string    `gorm:"column:label;type:text;not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"column:createdAt;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updatedAt;autoUpdateTime"`

	Gigs []Gig `gorm:"many2many:_GigToGigTag;foreignKey:ID;joinForeignKey:B;References:ID;joinReferences:A"`
}

func (GigTag) TableName() string { return "GigTag" }

// Gig maps to the "Gig" table

type Gig struct {
	ID            uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Title         string    `gorm:"column:title;type:text;not null"`
	Description   string    `gorm:"column:description;type:text;not null"`
	IsActive      bool      `gorm:"column:isActive;not null;default:true"`
	ViewCount     int       `gorm:"column:viewCount;not null;default:0"`
	AverageRating float64   `gorm:"column:averageRating;not null;default:0"`
	RatingCount   int       `gorm:"column:ratingCount;not null;default:0"`
	CategoryID    uuid.UUID `gorm:"column:categoryId;type:uuid;not null"`
	SellerID      uuid.UUID `gorm:"column:sellerId;type:uuid;not null"`
	CreatedAt     time.Time `gorm:"column:createdAt;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:updatedAt;autoUpdateTime"`

	Category Category     `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Seller   User         `gorm:"foreignKey:SellerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Images   []GigImage   `gorm:"foreignKey:GigID"`
	Packages []GigPackage `gorm:"foreignKey:GigID"`
	Tags     []GigTag     `gorm:"many2many:_GigToGigTag;foreignKey:ID;joinForeignKey:A;References:ID;joinReferences:B"`
}

func (Gig) TableName() string { return "Gig" }

// RegistrationToken maps to the "RegistrationToken" table

type RegistrationToken struct {
	Code      string     `gorm:"column:code;primaryKey;type:text"`
	Email     string     `gorm:"column:email;type:text;not null;uniqueIndex"`
	ExpiresAt time.Time  `gorm:"column:expiresAt;not null"`
	UserID    *uuid.UUID `gorm:"column:userId;type:uuid"`
	CreatedAt time.Time  `gorm:"column:createdAt;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updatedAt;autoUpdateTime"`

	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (RegistrationToken) TableName() string { return "RegistrationToken" }

// GigImage maps to the "GigImage" table

type GigImage struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	URL       string    `gorm:"column:url;type:text;not null"`
	Alt       *string   `gorm:"column:alt;type:text"`
	IsPrimary bool      `gorm:"column:isPrimary;not null;default:false"`
	SortOrder int       `gorm:"column:sortOrder;not null;default:0"`
	GigID     uuid.UUID `gorm:"column:gigId;type:uuid;not null"`
	CreatedAt time.Time `gorm:"column:createdAt;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updatedAt;autoUpdateTime"`

	Gig Gig `gorm:"foreignKey:GigID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (GigImage) TableName() string { return "GigImage" }

// GigPackage maps to the "GigPackage" table

type GigPackage struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Title        string    `gorm:"column:title;type:text;not null"`
	Description  string    `gorm:"column:description;type:text;not null"`
	Price        float64   `gorm:"column:price;not null"`
	DeliveryTime int       `gorm:"column:deliveryTime;not null"`
	Revisions    int       `gorm:"column:revisions;not null;default:1"`
	IsActive     bool      `gorm:"column:isActive;not null;default:true"`
	GigID        uuid.UUID `gorm:"column:gigId;type:uuid;not null"`
	CreatedAt    time.Time `gorm:"column:createdAt;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updatedAt;autoUpdateTime"`

	Gig      Gig                 `gorm:"foreignKey:GigID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Features []GigPackageFeature `gorm:"foreignKey:GigPackageID"`
}

func (GigPackage) TableName() string { return "GigPackage" }

// GigPackageFeature maps to the "GigPackageFeature" table

type GigPackageFeature struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Title        string    `gorm:"column:title;type:text;not null"`
	Description  *string   `gorm:"column:description;type:text"`
	Included     bool      `gorm:"column:included;not null;default:true"`
	GigPackageID uuid.UUID `gorm:"column:gigPackageId;type:uuid;not null"`
	CreatedAt    time.Time `gorm:"column:createdAt;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updatedAt;autoUpdateTime"`

	GigPackage GigPackage `gorm:"foreignKey:GigPackageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (GigPackageFeature) TableName() string { return "GigPackageFeature" }

// Category maps to the "Category" table

type Category struct {
	ID        uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Label     string     `gorm:"column:label;type:text;not null;uniqueIndex"`
	Slug      string     `gorm:"column:slug;type:text;not null;uniqueIndex"`
	IsActive  bool       `gorm:"column:isActive;not null;default:true"`
	SortOrder int        `gorm:"column:sortOrder;not null;default:0"`
	ParentID  *uuid.UUID `gorm:"column:parentId;type:uuid"`
	CreatedAt time.Time  `gorm:"column:createdAt;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updatedAt;autoUpdateTime"`

	Parent   *Category  `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Children []Category `gorm:"foreignKey:ParentID"`
	Gigs     []Gig      `gorm:"foreignKey:CategoryID"`
}

func (Category) TableName() string { return "Category" }

// Order maps to the "Order" table

type Order struct {
	ID            uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	OrderNumber   string     `gorm:"column:orderNumber;type:text;not null;uniqueIndex"`
	Price         float64    `gorm:"column:price;not null"`
	PaymentMethod string     `gorm:"column:paymentMethod;type:text;not null"`
	Status        string     `gorm:"column:status;type:text;not null;default:'PENDING'"`
	TransactionID *string    `gorm:"column:transactionId;type:text"`
	Requirements  *string    `gorm:"column:requirements;type:text"`
	PackageID     uuid.UUID  `gorm:"column:packageId;type:uuid;not null"`
	SellerID      uuid.UUID  `gorm:"column:sellerId;type:uuid;not null"`
	BuyerID       uuid.UUID  `gorm:"column:buyerId;type:uuid;not null"`
	CreatedAt     time.Time  `gorm:"column:createdAt;autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"column:updatedAt;autoUpdateTime"`
	CompletedAt   *time.Time `gorm:"column:completedAt"`
	DueDate       *time.Time `gorm:"column:dueDate"`

	Package GigPackage `gorm:"foreignKey:PackageID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Seller  User       `gorm:"foreignKey:SellerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Buyer   User       `gorm:"foreignKey:BuyerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Reviews []Review   `gorm:"foreignKey:OrderID"`
}

func (Order) TableName() string { return "Order" }

// Review maps to the "Review" table

type Review struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Title     string    `gorm:"column:title;type:text;not null"`
	Content   *string   `gorm:"column:content;type:text"`
	Rating    int       `gorm:"column:rating;not null;default:5"`
	IsPublic  bool      `gorm:"column:isPublic;not null;default:true"`
	AuthorID  uuid.UUID `gorm:"column:authorId;type:uuid;not null;index"`
	OrderID   uuid.UUID `gorm:"column:orderId;type:uuid;not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"column:createdAt;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updatedAt;autoUpdateTime"`

	Author User  `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Order  Order `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (Review) TableName() string { return "Review" }
