package service

import (
	"context"
	"fmt"

	"github.com/msmkdenis/yap-infokeeper/internal/model"

	apperr "github.com/msmkdenis/yap-infokeeper/pkg/apperror"
)

type CreditCardRepository interface {
	Insert(ctx context.Context, ownerID string, card model.CreditCard) error
}

type CreditCardUseCase struct {
	repository CreditCardRepository
}

func NewCreditCardService(repository CreditCardRepository) *CreditCardUseCase {
	return &CreditCardUseCase{repository: repository}
}

func (u *CreditCardUseCase) Save(ctx context.Context, ownerID string, card model.CreditCard) error {
	if err := u.repository.Insert(ctx, ownerID, card); err != nil {
		return fmt.Errorf("%s %w", apperr.Caller(), err)
	}

	return nil
}
