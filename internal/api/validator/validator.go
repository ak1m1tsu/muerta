package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func notBlank(fl validator.FieldLevel) bool {
	ok, err := regexp.Match(
		`^(?:[a-zA-Zа-яА-Я]|[a-zA-Zа-яА-Я]+\s)*[a-zA-Zа-яА-Я]$`,
		[]byte(fl.Field().String()),
	)
	if !ok || err != nil {
		return false
	}
	return true
}

func init() {
	validate = validator.New()
	validate.RegisterValidation("notblank", notBlank)
}

const (
	KeyErrResponses string = "validation"
	keyField        string = "field"
	keyTag          string = "tag"
	keyValue        string = "value"
)

type ValidationError struct {
	Field string      `json:"field"`
	Tag   string      `json:"tag"`
	Value interface{} `json:"value"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf(
		"`%v` with value `%v` doesn't satisfy the `%v` constraint",
		e.Field,
		e.Value,
		e.Tag,
	)
}

type ValidationErrors []ValidationError

func (errs ValidationErrors) Error() string {
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

func Validate(payload interface{}) ValidationErrors {
	if err := validate.Struct(payload); err != nil {
		var errors ValidationErrors
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(
				errors,
				ValidationError{
					Field: strings.ToLower(err.Field()),
					Tag:   err.Tag(),
					Value: err.Value(),
				},
			)
		}
		return errors
	}
	return nil
}
