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

func (h *UserRegister) PostRegisterUser(ctx context.Context, in *pb.PostUserRegisterRequest) (*pb.PostUserRegisterResponse, error) {
	user := model.User{
		ID:       in.Id,
		Login:    in.Login,
		Password: in.Password,
	}

	report, ok := h.validator.ValidateUser(user)
	if !ok {
		slog.Info("Unable to register user: invalid user request",
			slog.String("user_login", user.Login),
			slog.Any("violated_fields", report))
		return nil, processValidationError(report)
	}

	err := h.userService.Register(ctx, user)
	if errors.Is(err, model.ErrUserAlreadyExists) {
		slog.Info("Unable to register user: user already exists",
			slog.String("user_login", in.Login),
			slog.String("error", fmt.Errorf("%s %w", caller.CodeLine(), err).Error()))
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	}

	if err != nil {
		slog.Info("Unable to register user: internal error",
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

	return &pb.PostUserRegisterResponse{Token: token}, nil
}
