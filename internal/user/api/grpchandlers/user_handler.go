package grpchandlers

import (
	"context"
	"log/slog"
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
	pb "github.com/msmkdenis/yap-infokeeper/internal/proto/user"
	"github.com/msmkdenis/yap-infokeeper/pkg/jwtgen"
)

type UserService interface {
	Register(ctx context.Context, user model.User) error
	Login(ctx context.Context, user model.UserLoginRequest) (*model.User, error)
}

type UserRegister struct {
	userService UserService
	jwtManager  *jwtgen.JWTManager
	pb.UnimplementedUserServiceServer
	validator *model.UserRequestValidator
}

func NewUserRegister(service UserService, jwtManager *jwtgen.JWTManager) *UserRegister {
	validator := model.NewUserRequestValidator()
	return &UserRegister{
		userService: service,
		jwtManager:  jwtManager,
		validator:   validator,
	}
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
