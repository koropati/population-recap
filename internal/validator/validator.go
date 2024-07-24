package validator

import (
	"github.com/go-playground/validator/v10"
)

// Validator adalah struct yang menyimpan instance dari validator.
type Validator struct {
	validate *validator.Validate
}

// NewValidator membuat dan mengembalikan instance baru dari Validator.
func NewValidator() *Validator {
	return &Validator{
		validate: validator.New(),
	}
}

// Validate digunakan untuk memvalidasi struct berdasarkan tag yang didefinisikan pada struct tersebut.
func (v *Validator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}
