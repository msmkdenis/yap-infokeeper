package grpchandlers

import (
	"context"
	"log/slog"
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
	pb "github.com/msmkdenis/yap-infokeeper/internal/text_data/api/grpchandlers/proto"
	"github.com/msmkdenis/yap-infokeeper/internal/text_data/specification"
)

type TextDataService interface {
	Save(ctx context.Context, textData model.TextData) error
	Load(ctx context.Context, spec *specification.TextDataSpecification) ([]model.TextData, error)
}

type TextData struct {
	textDataService TextDataService
	pb.UnimplementedTextDataServiceServer
	validator *model.TextDataValidator
}

func NewTextData(textDataService TextDataService) *TextData {
	validator := model.NewTextDataValidator()
	return &TextData{
		textDataService: textDataService,
		validator:       validator,
	}
}

func processValidationError(report map[string][]string) error {
	st := status.New(codes.InvalidArgument, "invalid text data request")
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
