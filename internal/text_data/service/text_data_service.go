package service

import (
	"context"
	"fmt"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
	"github.com/msmkdenis/yap-infokeeper/internal/text_data/specification"
	apperr "github.com/msmkdenis/yap-infokeeper/pkg/apperror"
)

// mockgen --build_flags=--mod=mod -destination=internal/text_data/mocks/mock_text_data_repository.go -package=mocks github.com/msmkdenis/yap-infokeeper/internal/text_data/service TextDataRepository
type TextDataRepository interface {
	Insert(ctx context.Context, textData model.TextData) error
	SelectAll(ctx context.Context, spec *specification.TextDataSpecification) ([]model.TextData, error)
}

type TextDataUseCase struct {
	repository TextDataRepository
}

func NewTextDataService(repository TextDataRepository) *TextDataUseCase {
	return &TextDataUseCase{repository: repository}
}

func (u *TextDataUseCase) Save(ctx context.Context, textData model.TextData) error {
	if err := u.repository.Insert(ctx, textData); err != nil {
		return fmt.Errorf("%s %w", apperr.Caller(), err)
	}

	return nil
}

func (u *TextDataUseCase) Load(ctx context.Context, spec *specification.TextDataSpecification) ([]model.TextData, error) {
	textData, err := u.repository.SelectAll(ctx, spec)
	if err != nil {
		return nil, fmt.Errorf("%s %w", apperr.Caller(), err)
	}

	return textData, nil
}
