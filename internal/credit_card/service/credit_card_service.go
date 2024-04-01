package service

import (
	"context"
	"fmt"

	"github.com/msmkdenis/yap-infokeeper/internal/model"

	apperr "github.com/msmkdenis/yap-infokeeper/pkg/apperror"
)

// mockgen --build_flags=--mod=mod -destination=internal/credit_card/mocks/mock_credit_card_repository.go -package=mocks github.com/msmkdenis/yap-infokeeper/internal/credit_card/service CreditCardRepository
type CreditCardRepository interface {
	Insert(ctx context.Context, card model.CreditCard) error
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
