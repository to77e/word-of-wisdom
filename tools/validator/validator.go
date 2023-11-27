package validator

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func New() *validator.Validate {
	validate = validator.New(validator.WithRequiredStructEnabled())
	return validate
}

func Get() *validator.Validate {
	if validate != nil {
		return validate
	}

	slog.Warn("created new validator")
	return validator.New(validator.WithRequiredStructEnabled())
}
