package model

import (
	"fmt"
	"log/slog"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type CreditCard struct {
	ID        string    `db:"id" json:"id" validate:"uuid"`
	Number    string    `db:"number" json:"number" validate:"card_number"`
	OwnerID   string    `db:"owner_id" json:"owner_id" validate:"uuid"`
	OwnerName string    `db:"owner_name" json:"owner_name" validate:"owner"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CVVCode   string    `db:"cvv_code" json:"cvv_code" validate:"cvv"`
	PinCode   string    `db:"pin_code" json:"pin_code" validate:"pin"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Metadata  string    `db:"metadata" json:"metadata"`
}

type CreditCardRequestValidator struct {
	validator *validator.Validate
}

func NewCreditCardRequestValidator() *CreditCardRequestValidator {
	v := validator.New()

	err := v.RegisterValidation("cvv", cvvCode)
	if err != nil {
		slog.Error("unable to register validation cvv")
	}

	err = v.RegisterValidation("pin", pinCode)
	if err != nil {
		slog.Error("unable to register validation pin")
	}

	err = v.RegisterValidation("owner", owner)
	if err != nil {
		slog.Error("unable to register validation owner")
	}

	err = v.RegisterValidation("card_number", cardNumber)
	if err != nil {
		slog.Error("unable to register validation cardNumber")
	}

	return &CreditCardRequestValidator{validator: v}
}

func (v *CreditCardRequestValidator) ValidateCreditCard(request CreditCard) (map[string][]string, bool) {
	err := v.validator.Struct(request)
	report := make(map[string][]string)
	if err != nil {
		for _, validationErr := range err.(validator.ValidationErrors) { //nolint:errorlint
			switch validationErr.Tag() {
			case "card_number":
				report[validationErr.Field()] = append(report[validationErr.Field()], "must be valid credit card number")
			case "owner":
				report[validationErr.Field()] = append(report[validationErr.Field()], "must be valid owner")
			case "cvv":
				report[validationErr.Field()] = append(report[validationErr.Field()], "must be valid cvv")
			case "pin":
				report[validationErr.Field()] = append(report[validationErr.Field()], "must be valid pin")
			}
		}
		return report, false
	}
	return nil, true
}

func cardNumber(fl validator.FieldLevel) bool {
	block := strings.Split(fl.Field().String(), " ")
	if len(block) != 4 {
		return false
	}

	for _, block := range block {
		if len(block) != 4 {
			return false
		}
		for _, char := range block {
			if unicode.IsLetter(char) {
				return false
			}
		}
	}
	return true
}

func cvvCode(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) != 3 {
		return false
	}
	for _, char := range fl.Field().String() {
		if unicode.IsLetter(char) {
			return false
		}
	}
	return true
}

func pinCode(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) != 4 {
		return false
	}
	for _, char := range fl.Field().String() {
		if unicode.IsLetter(char) {
			return false
		}
	}
	return true
}

func owner(fl validator.FieldLevel) bool {
	if len(strings.Split(fl.Field().String(), " ")) != 2 {
		fmt.Println("owner", fl.Field().String())
		return false
	}
	return true
}
