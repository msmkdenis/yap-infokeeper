package pbclient

import (
	"context"

	pb "github.com/msmkdenis/yap-infokeeper/internal/proto/user"
)

type UserPBClient struct {
	userService pb.UserServiceClient
}

func NewUserPBClient(u pb.UserServiceClient) *UserPBClient {
	return &UserPBClient{
		userService: u,
	}
}

func (u *UserPBClient) LoginUser(login string, passwordHash []byte) (string, error) {
	req := &pb.PostUserLoginRequest{
		Login:    login,
		Password: passwordHash,
	}

	resp, err := u.userService.PostLoginUser(context.Background(), req)
	if err != nil {
		return "", err
	}

	return resp.Token, nil
}

func (u *UserPBClient) RegisterUser(id, login string, passwordHash []byte) (string, error) {
	req := &pb.PostUserRegisterRequest{
		Id:       id,
		Login:    login,
		Password: passwordHash,
	}

	resp, err := u.userService.PostRegisterUser(context.Background(), req)
	if err != nil {
		return "", err
	}

	return resp.Token, nil
}
