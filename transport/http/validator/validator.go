package validator

import (
	"errors"
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

type ApiError struct {
	PropertyName string `json:"propertyName"`
	ErrorMessage string `json:"errorMessage"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("%s: %s", e.PropertyName, e.ErrorMessage)
}

type ApiErrorResponse struct {
	Status     string
	StatusCode int
	Message    string
	Errors     []ApiError
}

func NewApiErrorResponse() ApiErrorResponse {
	return ApiErrorResponse{}
}

func (e ApiErrorResponse) Error() string {
	errs := ""
	for _, r := range e.Errors {
		errs = errs + fmt.Sprintf("%s: %s\n", r.PropertyName, r.ErrorMessage)
	}
	return errs
}

type CustomValidator struct {
	Validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	vl := validator.New()
	vl.RegisterValidation("date", ValidateDate)
	vl.RegisterValidation("time", ValidateTime)
	return &CustomValidator{Validator: vl}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func ValidatorMessageTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("Kolom %s diperlukan", fe.Field())
	case "email", "date", "time":
		return fmt.Sprintf("%s tidak valid", fe.Field())
	}
	return fe.Translate(nil) // default error
}

func ValidatorMessageFormat(err error) []ApiError {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ApiError, len(ve))
		for i, fe := range ve {
			out[i] = ApiError{fe.Field(), ValidatorMessageTag(fe)}
		}
		return out
	}
	return nil
}
