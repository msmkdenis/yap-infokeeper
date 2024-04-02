package grpchandlers

import (
	"context"
	"log/slog"
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/msmkdenis/yap-infokeeper/internal/credential/api/grpchandlers/proto"
	"github.com/msmkdenis/yap-infokeeper/internal/credential/specification"
	"github.com/msmkdenis/yap-infokeeper/internal/model"
)

// mockgen --build_flags=--mod=mod -destination=internal/credential/mocks/mock_credential_service.go -package=mocks github.com/msmkdenis/yap-infokeeper/internal/credential/api/grpchandlers CredentialService
type CredentialService interface {
	Save(ctx context.Context, credential model.Credential) error
	Load(ctx context.Context, spec *specification.CredentialSpecification) ([]model.Credential, error)
}

type Credential struct {
	credentialService CredentialService
	pb.UnimplementedCredentialServiceServer
	validator *model.CredentialValidator
}

func NewCredential(credentialService CredentialService) *Credential {
	validator := model.NewCredentialValidator()
	return &Credential{
		credentialService: credentialService,
		validator:         validator,
	}
}

func processValidationError(report map[string][]string) error {
	st := status.New(codes.InvalidArgument, "invalid credential request")
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
