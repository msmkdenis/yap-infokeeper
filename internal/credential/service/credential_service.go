package service

import (
	"context"
	"fmt"

	"github.com/msmkdenis/yap-infokeeper/internal/credential/specification"
	"github.com/msmkdenis/yap-infokeeper/internal/model"
	apperr "github.com/msmkdenis/yap-infokeeper/pkg/apperror"
)

type CredentialRepository interface {
	Insert(ctx context.Context, credential model.Credential) error
	SelectAll(ctx context.Context, spec *specification.CredentialSpecification) ([]model.Credential, error)
}

type CredentialUseCase struct {
	repository CredentialRepository
}

func NewCredentialService(repository CredentialRepository) *CredentialUseCase {
	return &CredentialUseCase{repository: repository}
}

func (u *CredentialUseCase) Save(ctx context.Context, credential model.Credential) error {
	if err := u.repository.Insert(ctx, credential); err != nil {
		return fmt.Errorf("%s %w", apperr.Caller(), err)
	}

	return nil
}

func (u *CredentialUseCase) Load(ctx context.Context, spec *specification.CredentialSpecification) ([]model.Credential, error) {
	credentials, err := u.repository.SelectAll(ctx, spec)
	if err != nil {
		return nil, fmt.Errorf("%s %w", apperr.Caller(), err)
	}

	return credentials, nil
}
