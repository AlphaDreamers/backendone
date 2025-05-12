package req

type SearchGigsRequest struct {
	Query      string   `query:"q" validate:"max=100"`
	Categories []string `query:"categories" validate:"max=3"`
	MinPrice   float64  `query:"minPrice" validate:"min=0"`
	MaxPrice   float64  `query:"maxPrice" validate:"gtfield=MinPrice"`
	SortBy     string   `query:"sortBy" validate:"oneof=price rating newest popular"`
	Page       int      `query:"page" validate:"min=1"`
	Limit      int      `query:"limit" validate:"min=1,max=50"`
}
