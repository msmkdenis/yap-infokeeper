package model

import (
	"github.com/go-playground/validator/v10"
)

type User struct {
	ID       string `db:"id" validate:"uuid"`
	Login    string `db:"login" validate:"email"`
	Password []byte `db:"password" validate:"required"`
}

type UserRequestValidator struct {
	validator *validator.Validate
}

func NewUserRequestValidator() *UserRequestValidator {
	validate := validator.New()
	return &UserRequestValidator{validator: validate}
}

func (v *UserRequestValidator) Validate(request User) (map[string][]string, bool) {
	err := v.validator.Struct(request)
	report := make(map[string][]string)
	if err != nil {
		for _, validationErr := range err.(validator.ValidationErrors) { //nolint:errorlint
			switch validationErr.Tag() {
			case "email":
				report[validationErr.Field()] = append(report[validationErr.Field()], "must be valid email")
			case "required":
				report[validationErr.Field()] = append(report[validationErr.Field()], "is required")
			case "uuid":
				report[validationErr.Field()] = append(report[validationErr.Field()], "must be valid uuid")
			}
		}
		return report, false
	}
	return nil, true
}
