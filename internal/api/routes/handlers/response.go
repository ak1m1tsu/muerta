package handlers

type HTTPSuccess struct {
	Success bool `json:"success" example:"true"`
	Data    Data `json:"data"`
}

type HTTPError struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"Not Found"`
}

type Data map[string]any
