package validator

import (
	"fmt"
	"strings"

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

func (er ErrorResponse) Error() string {
	return fmt.Sprintf("invalid field '%s' with tag '%s' value '%s'", er.Field, er.Tag, er.Value)
}

func (er ErrorResponse) MarshalZerologObject(e *zerolog.Event) {
	e.Str(keyField, er.Field).
		Str(keyTag, er.Tag).
		Str(keyValue, er.Value)
}

type ErrorResponses []ErrorResponse

func (errs ErrorResponses) Error() string {
	var buf strings.Builder
	for i, err := range errs {
		if len(errs)-1 == i {
			buf.WriteString(err.Error())
			break
		}
		buf.WriteString(fmt.Sprintf("%s, ", err))
	}
	return buf.String()
}

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
			errResp.Field = err.StructField()
			errResp.Tag = err.Tag()
			errResp.Value = err.Param()
			errors = append(errors, errResp)
		}
	}
	return errors
}
