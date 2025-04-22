package validation

import (
	"fmt"
	"regexp"
	"strings"
)


type Errors map[string]string

type Validator struct {
	errors Errors
}

func NewValidator() *Validator {
	return &Validator{
		errors: make(Errors),
	}
}

func (v *Validator) Valid() bool {
	return len(v.errors) == 0
}

func (v *Validator) Errors() Errors {
	return v.errors
}

func (v *Validator) Check(ok bool, field, message string) {
	if !ok {
		v.errors[field] = message
	}
}

func (v *Validator) Required(field, value string) {
	if strings.TrimSpace(value) == "" {
		v.errors[field] = "must not be empty"
	}
}

func (v *Validator) MinLength(field, value string, min int) {
	if len(value) < min {
		v.errors[field] = fmt.Sprintf("must be at least %d characters", min)
	}
}

func (v *Validator) MaxLength(field, value string, max int) {
    if len(value) > max {
        v.errors[field] = fmt.Sprintf("must not be more than %d characters", max)
    }
}

func (v *Validator) Email(field, value string) {
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    re := regexp.MustCompile(pattern)
    if !re.MatchString(value) {
        v.errors[field] = "must be a valid email address"
    }
}