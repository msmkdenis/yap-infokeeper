package service

import (
	"context"
	"fmt"

	apperr "github.com/msmkdenis/yap-infokeeper/pkg/apperror"

	"github.com/msmkdenis/yap-infokeeper/internal/user/model"
)

type UserRepository interface {
	Insert(ctx context.Context, user model.User) error
}

type UserUseCase struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserUseCase {
	return &UserUseCase{repository: repository}
}

func (u *UserUseCase) Register(ctx context.Context, user model.User) error {
	if err := u.repository.Insert(ctx, user); err != nil {
		return fmt.Errorf("%s %w", apperr.Caller(), err)
	}

	return nil
}
