package grpchandlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
	pb "github.com/msmkdenis/yap-infokeeper/internal/proto/user"
	"github.com/msmkdenis/yap-infokeeper/pkg/caller"
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
	if errors.Is(err, model.ErrInvalidPassword) {
		slog.Info("Unable to login user: invalid password",
			slog.String("user_login", in.Login),
			slog.String("error", fmt.Errorf("%s %w", caller.CodeLine(), err).Error()))
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	if errors.Is(err, model.ErrUserNotFound) {
		slog.Info("Unable to login user: user not found",
			slog.String("user_login", in.Login),
			slog.String("error", fmt.Errorf("%s %w", caller.CodeLine(), err).Error()))
		return nil, status.Error(codes.Unauthenticated, "user not found")
	}

	if err != nil {
		slog.Info("Unable to login user: internal error",
			slog.String("user_login", in.Login),
			slog.String("error", fmt.Errorf("%s %w", caller.CodeLine(), err).Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	token, err := h.jwtManager.BuildJWTString(user.ID)
	if err != nil {
		slog.Info("Unable to login user: failed to build token",
			slog.String("user_id", user.ID),
			slog.String("error", fmt.Errorf("%s %w", caller.CodeLine(), err).Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.PostUserLoginResponse{Token: token}, nil
}
