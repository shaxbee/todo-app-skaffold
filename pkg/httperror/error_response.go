package httperror

type ErrorResponse struct {
	Message string `json:"message"`
	Debug   string `json:"debug,omitempty"`
}
