package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// AlreadyExists is a custom error type for indicating that a record already exists.
type AlreadyExists struct {
	Message string
}

// Implement the Error() method for AlreadyExists to satisfy the error interface
func (ae *AlreadyExists) Error() string {
	return fmt.Sprintf("AlreadyExists: %s", ae.Message)
}

// DoesNotExist is a custom error type for indicating that a record doesn't exist.
type DoesNotExist struct {
	Message string
}

// Implement the Error() method for DoesNotExist to satisfy the error interface
func (ae *DoesNotExist) Error() string {
	return fmt.Sprintf("DoesNotExist: %s", ae.Message)
}

// Invalid is a custom error type for indicating that a request is invalid.
type InvalidRequest struct {
	Message string
}

// Implement the Error() method for Invalid to satisfy the error interface
func (ae *InvalidRequest) Error() string {
	return fmt.Sprintf("InvalidRequest: %s", ae.Message)
}

type ValidationError struct {
	Field string
	Err   error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Err.Error())
}

func getValidationErrors(err error) []ValidationError {
	validationErrors := []ValidationError{}

	for _, fieldErr := range err.(validator.ValidationErrors) {
		validationErrors = append(validationErrors, ValidationError{
			Field: fieldErr.StructField(),
			Err:   fieldErr,
		})
	}

	return validationErrors
}
