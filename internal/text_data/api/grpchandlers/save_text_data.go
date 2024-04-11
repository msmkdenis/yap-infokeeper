package grpchandlers

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-infokeeper/internal/interceptors"
	"github.com/msmkdenis/yap-infokeeper/internal/model"
	pb "github.com/msmkdenis/yap-infokeeper/internal/proto/text_data"
)

func (h *TextData) PostSaveTextData(ctx context.Context, in *pb.PostTextDataRequest) (*pb.PostTextDataResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDContextKey("userID")).(string)
	if !ok {
		slog.Error("Internal server error", slog.String("error", "unable to get userID from context"))
		return nil, status.Error(codes.Internal, "internal error")
	}

	texData := model.TextData{
		ID:       in.Uuid,
		OwnerID:  userID,
		Data:     in.Data,
		Metadata: in.Metadata,
	}

	report, ok := h.validator.ValidateTextData(texData)
	if !ok {
		return nil, processValidationError(report)
	}

	err := h.textDataService.Save(ctx, texData)
	if err != nil {
		slog.Info("Unable to save text data", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error while saving text data")
	}

	return &pb.PostTextDataResponse{}, nil
}
