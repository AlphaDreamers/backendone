package response

type Response struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"message"`
	Data    any    `json:"body,omitempty"`
}
