package util

import (
	"github.com/go-playground/validator/v10"
)

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

func SuppotedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD:
		return true
	}
	return false
}

var ValidCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return SuppotedCurrency(currency)
	}
	return false
}
