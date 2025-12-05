package model

// ErrorResponse defines the standardised JSON structure for API errors.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
