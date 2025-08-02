package validator

import (
	"net/mail"
)

type Validator struct {
	ErrVal map[string]string
}

func NewValidator() *Validator {
	errs := make(map[string]string)
	return &Validator{ErrVal: errs}
}

func (v *Validator) Add(key, message string) {
	if _, exist := v.ErrVal[key]; !exist {
		v.ErrVal[key] = message
	}
}

func (v *Validator) Check(ok bool, key, mesage string) {
	if !ok {
		v.Add(key, mesage)
	}
}

func (v *Validator) Valid() bool {
	return len(v.ErrVal) == 0
}

func (v *Validator) Error() string {
	return "failed validation"
}

func EmailCheck(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
