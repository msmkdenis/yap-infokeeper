package grpchandlers

import (
	"context"
	"log/slog"
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/msmkdenis/yap-infokeeper/internal/credit_card/api/grpchandlers/proto"
	"github.com/msmkdenis/yap-infokeeper/internal/model"
)

// mockgen --build_flags=--mod=mod -destination=internal/credit_card/mocks/mock_credit_card_service.go -package=mocks github.com/msmkdenis/yap-infokeeper/internal/credit_card/api/grpchandlers CreditCardService
type CreditCardService interface {
	Save(ctx context.Context, card model.CreditCard) error
	SelectByOwnerIDCardNumber(ctx context.Context, ownerID, number string) (*model.CreditCard, error)
	SelectAllByOwnerID(ctx context.Context, ownerID string) ([]model.CreditCard, error)
}

type CreditCard struct {
	creditCardService CreditCardService
	pb.UnimplementedCreditCardServiceServer
	validator *model.CreditCardRequestValidator
}

func NewCreditCard(creditCardService CreditCardService) *CreditCard {
	validator := model.NewCreditCardRequestValidator()
	return &CreditCard{
		creditCardService: creditCardService,
		validator:         validator,
	}
}

func processValidationError(report map[string][]string) error {
	st := status.New(codes.InvalidArgument, "invalid credit card request")
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
