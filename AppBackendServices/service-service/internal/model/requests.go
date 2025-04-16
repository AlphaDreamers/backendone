package model

// ServiceRequest represents a service creation/update request
type ServiceRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Category    string  `json:"category" validate:"required"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Status      string  `json:"status" validate:"required,oneof=active inactive pending"`
}

// PaginationRequest represents a pagination request
type PaginationRequest struct {
	Page     int `query:"page" validate:"min=1"`
	PageSize int `query:"page_size" validate:"min=1,max=100"`
}

// FilterRequest represents a filter request
type FilterRequest struct {
	Category  string `query:"category"`
	Status    string `query:"status"`
	SortBy    string `query:"sort_by"`
	SortOrder string `query:"sort_order" validate:"oneof=asc desc"`
}

//// TransactionFilterRequest represents a transaction filter request
//type TransactionFilterRequest struct {
//	Type      string `query:"type" validate:"oneof=fiat crypto all"`
//	StartDate string `query:"start_date"`
//	EndDate   string `query:"end_date"`
//	Status    string `query:"status"`
//}
//
//// ChatMessageRequest represents a chat message request
//type ChatMessageRequest struct {
//	Message string `json:"message" validate:"required"`
//}
