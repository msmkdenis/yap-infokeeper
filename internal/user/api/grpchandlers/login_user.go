package grpchandlers

import (
	"context"
	"errors"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
	pb "github.com/msmkdenis/yap-infokeeper/internal/proto/user"
	apperr "github.com/msmkdenis/yap-infokeeper/pkg/apperror"
)

func (h *UserRegister) PostLoginUser(ctx context.Context, in *pb.PostUserLoginRequest) (*pb.PostUserLoginResponse, error) {
	userLoginRequest := model.UserLoginRequest{
		Login:    in.Login,
		Password: in.Password,
	}

	report, ok := h.validator.ValidateUserLoginRequest(userLoginRequest)
	if !ok {
		return nil, processValidationError(report)
	}

	user, err := h.userService.Login(ctx, userLoginRequest)
	if errors.Is(err, apperr.ErrInvalidPassword) {
		slog.Info("Invalid password", slog.String("with login", in.Login))
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	if errors.Is(err, apperr.ErrUserNotFound) {
		slog.Info("User not found", slog.String("with login", in.Login))
		return nil, status.Error(codes.Unauthenticated, "user not found")
	}

	if err != nil {
		slog.Info("Unable to login user", slog.String("with login", in.Login))
		return nil, status.Error(codes.Internal, "internal error")
	}

	token, err := h.jwtManager.BuildJWTString(user.ID)
	if err != nil {
		slog.Error("failed to build token", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.PostUserLoginResponse{Token: token}, nil
}
