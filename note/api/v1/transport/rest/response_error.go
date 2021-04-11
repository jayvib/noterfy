package rest

// ResponseError is the container to any error response.
type ResponseError struct {
	Message string `json:"message,omitempty"`
}
