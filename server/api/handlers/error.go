package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field string
	Err   error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Err.Error())
}

func getValidationErrors(err error) string {
	validationErrors := []ValidationError{}

	for _, fieldErr := range err.(validator.ValidationErrors) {
		validationErrors = append(validationErrors, ValidationError{
			Field: fieldErr.StructField(),
			Err:   fieldErr,
		})
	}

	return fmt.Sprint(validationErrors)
}
