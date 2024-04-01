package service

import (
	"context"
	"fmt"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
	apperr "github.com/msmkdenis/yap-infokeeper/pkg/apperror"
)

// mockgen --build_flags=--mod=mod -destination=internal/credential/mocks/mock_credential_repository.go -package=mocks github.com/msmkdenis/yap-infokeeper/internal/credential/service CredentialRepository
type CredentialRepository interface {
	Insert(ctx context.Context, credential model.Credential) error
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
