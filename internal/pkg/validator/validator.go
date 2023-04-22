package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

var validate = validator.New()

const (
	KeyErrResponses string = "errorResponses"
	keyField        string = "field"
	keyTag          string = "tag"
	keyValue        string = "value"
)

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

func (er ErrorResponse) MarshalZerologObject(e *zerolog.Event) {
	e.Str(keyField, er.Field).
		Str(keyTag, er.Tag).
		Str(keyValue, er.Value)
}

type ErrorResponses []ErrorResponse

func (ers ErrorResponses) MarshalZerologArray(a *zerolog.Array) {
	for _, er := range ers {
		a.Object(er)
	}
}

func Validate(payload interface{}) ErrorResponses {
	var errors ErrorResponses
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var errResp ErrorResponse
			errResp.Field = err.StructNamespace()
			errResp.Tag = err.Tag()
			errResp.Value = err.Param()
			errors = append(errors, errResp)
		}
	}
	return errors
}
