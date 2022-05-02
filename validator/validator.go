package validator

import "net/mail"

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) AddError(key string, val string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = val
	}
}

func (v *Validator) IsEmail(email string) bool {
	if _, err := mail.ParseAddress(email); err != nil {
		return false
	}
	return true
}

func (v *Validator) Check(ok bool, key string, val string) {
	if !ok {
		v.AddError(key, val)
	}
}

func (v *Validator) Valid() bool {
	if len(v.Errors) == 0 {
		return true
	}
	return false
}
