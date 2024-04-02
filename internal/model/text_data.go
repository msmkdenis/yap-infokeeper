package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type TextData struct {
	ID        string    `db:"id" validate:"uuid"`
	OwnerID   string    `db:"owner_id" validate:"uuid"`
	Data      string    `db:"data" validate:"required"`
	CreatedAt time.Time `db:"created_at"`
	Metadata  string    `db:"metadata"`
}

type TextDataValidator struct {
	validator *validator.Validate
}

func NewTextDataValidator() *TextDataValidator {
	v := validator.New()
	return &TextDataValidator{validator: v}
}

func (v *TextDataValidator) ValidateTextData(request TextData) (map[string][]string, bool) {
	err := v.validator.Struct(request)
	report := make(map[string][]string)
	if err != nil {
		for _, validationErr := range err.(validator.ValidationErrors) { //nolint:errorlint
			switch validationErr.Tag() {
			case UUID:
				report[validationErr.Field()] = append(report[validationErr.Field()], "must be valid uuid")
			case REQUIRED:
				report[validationErr.Field()] = append(report[validationErr.Field()], "must be not empty")
			}
		}
		return report, false
	}
	return nil, true
}
