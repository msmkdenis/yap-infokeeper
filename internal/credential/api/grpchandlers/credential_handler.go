package grpchandlers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-infokeeper/internal/credential/specification"
	"github.com/msmkdenis/yap-infokeeper/internal/model"
	pb "github.com/msmkdenis/yap-infokeeper/internal/proto/credential"
	"github.com/msmkdenis/yap-infokeeper/pkg/caller"
)

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
		slog.Error("Internal error: failed to set details",
			slog.String("error", fmt.Errorf("%s %w", caller.CodeLine(), err).Error()))
		return status.Error(codes.Internal, "internal error")
	}
	return st.Err()
}
