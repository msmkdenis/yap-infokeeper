package grpchandlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/msmkdenis/yap-infokeeper/internal/user/api/grpchandlers/proto"
	"github.com/msmkdenis/yap-infokeeper/internal/user/model"
)

type UserService interface {
	Register(ctx context.Context, user model.User) error
}

type UserRegister struct {
	userService UserService
	pb.UnimplementedUserServiceServer
}

func NewUserRegister(service UserService) *UserRegister {
	return &UserRegister{userService: service}
}

func (h *UserRegister) PostRegisterUser(ctx context.Context, in *pb.PostUserRegisterRequest) (*pb.PostUserRegisterResponse, error) {
	user := model.User{
		ID:       in.Id,
		Login:    in.Login,
		Password: in.Password,
	}

	err := h.userService.Register(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.PostUserRegisterResponse{Token: "token"}, nil
}
