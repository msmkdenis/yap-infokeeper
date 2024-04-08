package service

import (
	"context"
	"fmt"

	"github.com/msmkdenis/yap-infokeeper/internal/credit_card/specification"
	"github.com/msmkdenis/yap-infokeeper/internal/model"
	apperr "github.com/msmkdenis/yap-infokeeper/pkg/apperror"
)

type CreditCardRepository interface {
	Insert(ctx context.Context, card model.CreditCard) error
	SelectByOwnerIDCardNumber(ctx context.Context, ownerID, number string) (*model.CreditCard, error)
	SelectAllByOwnerID(ctx context.Context, ownerID string) ([]model.CreditCard, error)
	SelectAll(ctx context.Context, spec *specification.CreditCardSpecification) ([]model.CreditCard, error)
}

type CreditCardUseCase struct {
	repository CreditCardRepository
}

func NewCreditCardService(repository CreditCardRepository) *CreditCardUseCase {
	return &CreditCardUseCase{repository: repository}
}

func (u *CreditCardUseCase) Save(ctx context.Context, card model.CreditCard) error {
	if err := u.repository.Insert(ctx, card); err != nil {
		return fmt.Errorf("%s %w", apperr.Caller(), err)
	}

	return nil
}

func (u *CreditCardUseCase) Load(ctx context.Context, spec *specification.CreditCardSpecification) ([]model.CreditCard, error) {
	cards, err := u.repository.SelectAll(ctx, spec)
	if err != nil {
		return nil, fmt.Errorf("%s %w", apperr.Caller(), err)
	}

	return cards, nil
}
