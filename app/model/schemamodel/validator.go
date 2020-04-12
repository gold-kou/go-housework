package schemamodel

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const (
	// MinUserNameLength the minimum length of user name
	MinUserNameLength = 1
	//MaxUserNameLength the max length of user name
	MaxUserNameLength = 100
	// MinPasswordLength the minimum length of password
	MinPasswordLength = 6
	// MaxPasswordLength the max length of password
	MaxPasswordLength = 20
	// MinFamilyNameLength the minimum length of family name
	MinFamilyNameLength = 1
	// MaxFamilyNameLength the max length of family
	MaxFamilyNameLength = 100
)

// RequestSchemaModel represents interface of request schema
type RequestSchemaModel interface {
	ValidateParam() error
}

// ValidateParam validation
func (cu *RequestCreateUser) ValidateParam() error {
	var fieldRules []*validation.FieldRules
	fieldRules = append(fieldRules, validation.Field(&cu.Email, validation.Required, is.Email))
	fieldRules = append(fieldRules, validation.Field(&cu.UserName, validation.Required,
		validation.Length(MinUserNameLength, MaxUserNameLength)))
	fieldRules = append(fieldRules, validation.Field(&cu.Password, validation.Required,
		validation.Length(MinPasswordLength, MaxPasswordLength)))
	return validation.ValidateStruct(cu, fieldRules...)
}

// ValidateParam validation
func (cf *RequestCreateFamily) ValidateParam() error {
	var fieldRules []*validation.FieldRules
	fieldRules = append(fieldRules, validation.Field(&cf.FamilyName, validation.Required,
		validation.Length(MinFamilyNameLength, MaxFamilyNameLength)))
	return validation.ValidateStruct(cf, fieldRules...)
}

// ValidateParam validation
func (uf *RequestUpdateFamily) ValidateParam() error {
	var fieldRules []*validation.FieldRules
	fieldRules = append(fieldRules, validation.Field(&uf.FamilyName, validation.Required,
		validation.Length(MinFamilyNameLength, MaxFamilyNameLength)))
	return validation.ValidateStruct(uf, fieldRules...)
}

// ValidateParam validation
func (rfm *RequestRegisterFamilyMember) ValidateParam() error {
	var fieldRules []*validation.FieldRules
	fieldRules = append(fieldRules, validation.Field(&rfm.MemberName, validation.Required))
	return validation.ValidateStruct(rfm, fieldRules...)
}

// ValidateParam validation
func (rct *RequestCreateTask) ValidateParam() error {
	var fieldRules []*validation.FieldRules
	fieldRules = append(fieldRules, validation.Field(&rct.TaskName, validation.Required))
	fieldRules = append(fieldRules, validation.Field(&rct.MemberName, validation.Required))
	fieldRules = append(fieldRules, validation.Field(&rct.Status, validation.Required, validation.In("todo", "done")))
	fieldRules = append(fieldRules, validation.Field(&rct.Date, validation.Required, validation.Date("2006-01-02")))
	return validation.ValidateStruct(rct, fieldRules...)
}

// ValidateParam validation
func (rut *RequestUpdateTask) ValidateParam() error {
	var fieldRules []*validation.FieldRules
	fieldRules = append(fieldRules, validation.Field(&rut.Task))
	return validation.ValidateStruct(rut, fieldRules...)
}

/*
func (t *Task) Validate() error {
	var fieldRules []*validation.FieldRules
	if t != nil {
		fieldRules = append(fieldRules, validation.Field(&t.TaskId, validation.Required))
		fieldRules = append(fieldRules, validation.Field(&t.TaskName, validation.Required))
		fieldRules = append(fieldRules, validation.Field(&t.MemberName, validation.Required))
		fieldRules = append(fieldRules, validation.Field(&t.Status, validation.Required, validation.In("todo", "done")))
		fieldRules = append(fieldRules, validation.Field(&t.Date, validation.Required, validation.Date("2006-01-02")))
	}
	return validation.ValidateStruct(t, fieldRules...)
}*/
