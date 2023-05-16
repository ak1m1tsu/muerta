package handlers

type HTTPSuccess struct {
	Success bool `json:"success" example:"true"`
	Data    Data `json:"data,omitempty"`
}

type HTTPError struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error,omitempty" example:"Not Found"`
}

type Data map[string]any
