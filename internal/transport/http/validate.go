package httptransport

import (
	"github.com/0xOnah/bank/internal/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return util.SuppotedCurrency(currency)
	}
	return false
}
