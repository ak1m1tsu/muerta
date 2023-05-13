package handlers

type apiResponse struct {
	Success bool   `json:"success"`
	Data    Data   `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (r *apiResponse) WithSuccess() *apiResponse {
	r.Success = true
	return r
}

func (r *apiResponse) WithData(data Data) *apiResponse {
	r.Data = data
	return r
}

type Data map[string]any

func SuccessResponse() *apiResponse {
	return &apiResponse{
		Success: true,
	}
}

func ErrorResponse(err error) *apiResponse {
	return &apiResponse{
		Success: false,
		Error:   err.Error(),
	}
}
