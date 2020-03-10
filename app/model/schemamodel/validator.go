package schemamodel

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const (
	MIN_USER_NAME_LENGTH = 1
	MAX_USER_NAME_LENGTH = 100
	MIN_PASSWORD_LENGTH  = 6
	MAX_PASSWORD_LENGTH  = 20
)

func (user *RequestCreateUser) ValidateParam() error {
	var fieldRules []*validation.FieldRules
	fieldRules = append(fieldRules, validation.Field(&user.Email, validation.Required, is.Email))
	fieldRules = append(fieldRules, validation.Field(&user.UserName, validation.Required,
		validation.Length(MIN_USER_NAME_LENGTH, MAX_USER_NAME_LENGTH)))
	fieldRules = append(fieldRules, validation.Field(&user.Password, validation.Required,
		validation.Length(MIN_PASSWORD_LENGTH, MAX_PASSWORD_LENGTH)))
	return validation.ValidateStruct(user, fieldRules...)
}
