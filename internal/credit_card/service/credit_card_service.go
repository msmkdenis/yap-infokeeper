package service

import (
	"context"
	"fmt"

	"github.com/msmkdenis/yap-infokeeper/internal/credit_card/specification"
	"github.com/msmkdenis/yap-infokeeper/internal/model"
	"github.com/msmkdenis/yap-infokeeper/pkg/caller"
)

type CreditCardRepository interface {
	Insert(ctx context.Context, card model.CreditCard) error
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
		return fmt.Errorf("%s %w", caller.CodeLine(), err)
	}

	return nil
}

func (u *CreditCardUseCase) Load(ctx context.Context, spec *specification.CreditCardSpecification) ([]model.CreditCard, error) {
	cards, err := u.repository.SelectAll(ctx, spec)
	if err != nil {
		return nil, fmt.Errorf("%s %w", caller.CodeLine(), err)
	}

	return cards, nil
}
