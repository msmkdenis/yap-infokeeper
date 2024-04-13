package service

import (
	"bytes"
	"context"
	"fmt"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
	"github.com/msmkdenis/yap-infokeeper/pkg/caller"
)

type UserRepository interface {
	Insert(ctx context.Context, user model.User) error
	SelectByLogin(ctx context.Context, login string) (*model.User, error)
}

type UserUseCase struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserUseCase {
	return &UserUseCase{repository: repository}
}

func (u *UserUseCase) Register(ctx context.Context, user model.User) error {
	if err := u.repository.Insert(ctx, user); err != nil {
		return fmt.Errorf("%s %w", caller.CodeLine(), err)
	}

	return nil
}

func (u *UserUseCase) Login(ctx context.Context, userLoginRequest model.UserLoginRequest) (*model.User, error) {
	user, err := u.repository.SelectByLogin(ctx, userLoginRequest.Login)
	if err != nil {
		return nil, fmt.Errorf("%s %w", caller.CodeLine(), err)
	}

	if !bytes.Equal(userLoginRequest.Password, user.Password) {
		return nil, fmt.Errorf("%s %w", caller.CodeLine(), model.ErrInvalidPassword)
	}

	return user, nil
}
