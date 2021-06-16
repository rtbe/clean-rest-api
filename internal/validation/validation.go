// Package validation provides validation for incoming data.
package validation

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// validate holds validation related settings.
var validate *validator.Validate

// translator is a cache for supported translation locale.
var translator ut.Translator

func init() {
	// Initialize a validator.
	validate = validator.New()

	// Create a english translator.
	translator, _ = ut.New(en.New(), en.New()).GetTranslator("en")

	// Register english translations as default translations for validation error messages.
	en_translations.RegisterDefaultTranslations(validate, translator)
}

func Check(v interface{}) error {
	if err := validate.Struct(v); err != nil {

		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		var fieldErrors FieldErrors
		for _, validationError := range validationErrors {
			fieldError := FieldError{
				Field: validationError.Field(),
				Error: validationError.Error(),
			}
			fieldErrors = append(fieldErrors, fieldError)
		}
		return fieldErrors
	}

	return nil
}
