package httptransport

import (
	"github.com/go-playground/validator/v10"
	"github.com/onahvictor/bank/internal/util"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.SuppotedCurrency(currency)
	}
	return false
}
