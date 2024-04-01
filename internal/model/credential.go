package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Credential struct {
	ID        string    `db:"id" validate:"uuid"`
	OwnerID   string    `db:"owner_id" validate:"uuid"`
	Login     string    `db:"login" validate:"required"`
	Password  string    `db:"password" validate:"required"`
	CreatedAt time.Time `db:"created_at"`
	Metadata  string    `db:"metadata"`
}

type CredentialValidator struct {
	validator *validator.Validate
}

func NewCredentialValidator() *CredentialValidator {
	v := validator.New()
	return &CredentialValidator{validator: v}
}

func (v *CredentialValidator) ValidateCredential(request Credential) (map[string][]string, bool) {
	err := v.validator.Struct(request)
	report := make(map[string][]string)
	if err != nil {
		for _, validationErr := range err.(validator.ValidationErrors) { //nolint:errorlint
			switch validationErr.Tag() {
			case "uuid":
				report[validationErr.Field()] = append(report[validationErr.Field()], "must be valid uuid")
			case "required":
				report[validationErr.Field()] = append(report[validationErr.Field()], "must be not empty")
			}
		}
		return report, false
	}
	return nil, true
}
