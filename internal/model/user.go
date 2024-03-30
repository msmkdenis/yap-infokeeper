package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

const (
	EMAIL    = "email"
	REQUIRED = "required"
	UUID     = "uuid"
)

type User struct {
	ID        string    `db:"id" validate:"uuid"`
	Login     string    `db:"login" validate:"email"`
	Password  []byte    `db:"password" validate:"required"`
	CreatedAt time.Time `db:"created_at"`
}

type UserLoginRequest struct {
	Login    string `validate:"email"`
	Password []byte `validate:"required"`
}

type UserRequestValidator struct {
	validator *validator.Validate
}

func NewUserRequestValidator() *UserRequestValidator {
	validate := validator.New()
	return &UserRequestValidator{validator: validate}
}

func (v *UserRequestValidator) ValidateUser(request User) (map[string][]string, bool) {
	err := v.validator.Struct(request)
	report := make(map[string][]string)
	if err != nil {
		for _, validationErr := range err.(validator.ValidationErrors) { //nolint:errorlint
			switch validationErr.Tag() {
			case EMAIL:
				report[validationErr.Field()] = append(report[validationErr.Field()], "must be valid email")
			case REQUIRED:
				report[validationErr.Field()] = append(report[validationErr.Field()], "is required")
			case UUID:
				report[validationErr.Field()] = append(report[validationErr.Field()], "must be valid uuid")
			}
		}
		return report, false
	}
	return nil, true
}

func (v *UserRequestValidator) ValidateUserLoginRequest(request UserLoginRequest) (map[string][]string, bool) {
	err := v.validator.Struct(request)
	report := make(map[string][]string)
	if err != nil {
		for _, validationErr := range err.(validator.ValidationErrors) { //nolint:errorlint
			switch validationErr.Tag() {
			case EMAIL:
				report[validationErr.Field()] = append(report[validationErr.Field()], "must be valid email")
			case REQUIRED:
				report[validationErr.Field()] = append(report[validationErr.Field()], "is required")
			}
		}
		return report, false
	}
	return nil, true
}
