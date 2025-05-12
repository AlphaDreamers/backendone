package req

type CreateGigRequest struct {
	Title       string              `json:"title" validate:"required,min=10,max=100"`
	Description string              `json:"description" validate:"required,min=50,max=1000"`
	CategoryID  string              `json:"categoryName" validate:"required,uuid"`
	SellerID    string              `json:"sellerId" validate:"required,uuid"`
	Tags        []string            `json:"tags" validate:"max=5"`
	Packages    []GigPackageRequest `json:"packages" validate:"required,min=1,max=3"`
	Images      []GigImageRequest   `json:"images" validate:"max=5"`
}

type GigPackageRequest struct {
	Title        string           `json:"title" validate:"required,min=5,max=50"`
	Description  string           `json:"description" validate:"required,min=20,max=200"`
	Price        float64          `json:"price" validate:"required,min=5,max=10000"`
	DeliveryDays int              `json:"deliveryDays" validate:"required,min=1,max=90"`
	Features     []FeatureRequest `json:"features" validate:"max=10"`
}
type FeatureRequest struct {
	Title       string `json:"title" validate:"required,min=5,max=50"`
	Description string `json:"description" validate:"required,min=20,max=200"`
	Included    bool   `json:"included" validate:"required"`
}

type GigImageRequest struct {
	URL       string `json:"url" validate:"required,url"`
	AltText   string `json:"altText" validate:"max=100"`
	IsPrimary bool   `json:"isPrimary"`
}
