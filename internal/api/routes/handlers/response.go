package handlers

type apiResponse struct {
	Success bool `json:"success"`
	Data    Data `json:"data,omitempty" swaggertype:"object,string" example:"key:value,key2:value"`
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

type errorResponse struct {
	apiResponse
	Error string `json:"error,omitempty" example:"error message"`
}

func ErrorResponse(err error) *errorResponse {
	return &errorResponse{
		apiResponse: apiResponse{Success: false},
		Error:       err.Error(),
	}
}
