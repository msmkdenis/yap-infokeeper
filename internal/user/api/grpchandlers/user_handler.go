package grpchandlers

import (
	"context"
	"log/slog"
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/msmkdenis/yap-infokeeper/internal/user/api/grpchandlers/proto"
	"github.com/msmkdenis/yap-infokeeper/internal/user/model"
	"github.com/msmkdenis/yap-infokeeper/pkg/jwtgen"
)

type UserService interface {
	Register(ctx context.Context, user model.User) error
}

type UserRegister struct {
	userService UserService
	jwtManager  *jwtgen.JWTManager
	pb.UnimplementedUserServiceServer
}

func NewUserRegister(service UserService, jwtManager *jwtgen.JWTManager) *UserRegister {
	return &UserRegister{
		userService: service,
		jwtManager:  jwtManager,
	}
}

func (h *UserRegister) PostRegisterUser(ctx context.Context, in *pb.PostUserRegisterRequest) (*pb.PostUserRegisterResponse, error) {
	user := model.User{
		ID:       in.Id,
		Login:    in.Login,
		Password: in.Password,
	}

	validator := model.NewUserRequestValidator()
	report, ok := validator.Validate(user)
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

func processValidationError(report map[string][]string) error {
	st := status.New(codes.InvalidArgument, "invalid user request")
	details := make([]*errdetails.BadRequest_FieldViolation, 0, len(report))
	for field, messages := range report {
		var description strings.Builder
		for _, message := range messages {
			description.WriteString(message)
		}
		details = append(details, &errdetails.BadRequest_FieldViolation{
			Field:       field,
			Description: description.String(),
		})
	}
	br := &errdetails.BadRequest{}
	br.FieldViolations = append(br.FieldViolations, details...)
	st, err := st.WithDetails(br)
	if err != nil {
		slog.Error("failed to set details", slog.String("error", err.Error()))
		return status.Error(codes.Internal, "internal error")
	}
	slog.Info("validation error", slog.String("error", st.Err().Error()))
	return st.Err()
}
