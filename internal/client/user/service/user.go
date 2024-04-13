package service

import "strings"

type UserProcessor interface {
	RegisterUser(id, login string, passwordHash []byte) (string, error)
	LoginUser(login string, passwordHash []byte) (string, error)
}

type UserService struct {
	userClient UserProcessor
}

func NewUserService(u UserProcessor) *UserService {
	return &UserService{
		userClient: u,
	}
}

func (u *UserService) RegisterUser(data string) string {
	args := strings.Fields(data)
	token, err := u.userClient.RegisterUser(args[0], args[1], []byte(args[2]))
	if err != nil {
		return err.Error()
	}

	return token
}

func (u *UserService) LoginUser(data string) string {
	args := strings.Fields(data)
	token, err := u.userClient.LoginUser(args[0], []byte(args[1]))
	if err != nil {
		return err.Error()
	}

	return token
}
