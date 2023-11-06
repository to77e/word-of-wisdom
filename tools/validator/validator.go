package validator

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func New() *validator.Validate {
	validate = validator.New(validator.WithRequiredStructEnabled())
	return validate
}

func Get() *validator.Validate {
	return validate
}
