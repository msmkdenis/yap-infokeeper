package grpchandlers

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
	pb "github.com/msmkdenis/yap-infokeeper/internal/proto/user"
)

func (h *UserRegister) PostRegisterUser(ctx context.Context, in *pb.PostUserRegisterRequest) (*pb.PostUserRegisterResponse, error) {
	user := model.User{
		ID:       in.Id,
		Login:    in.Login,
		Password: in.Password,
	}

	report, ok := h.validator.ValidateUser(user)
	if !ok {
		return nil, processValidationError(report)
	}

	err := h.userService.Register(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	token, err := h.jwtManager.BuildJWTString(user.ID)
	if err != nil {
		slog.Error("failed to build token", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.PostUserRegisterResponse{Token: token}, nil
}
