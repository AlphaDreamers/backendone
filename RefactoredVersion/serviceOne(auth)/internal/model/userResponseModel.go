package model

// DashboardRespons
type DashboardResponse struct {
	UserInfo   *User         `json:"user_info"`
	Services   []ServicePost `json:"services"`
	Pagination struct {
		Page       int   `json:"page"`
		PageSize   int   `json:"page_size"`
		Total      int64 `json:"total_items"`
		TotalPages int   `json:"total_pages"`
	} `json:"pagination"`
}
