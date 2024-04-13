package grpchandlers

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-infokeeper/internal/interceptors"
	"github.com/msmkdenis/yap-infokeeper/internal/model"
	pb "github.com/msmkdenis/yap-infokeeper/internal/proto/text_data"
	"github.com/msmkdenis/yap-infokeeper/pkg/caller"
)

func (h *TextData) PostSaveTextData(ctx context.Context, in *pb.PostTextDataRequest) (*pb.PostTextDataResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDContextKey("userID")).(string)
	if !ok {
		slog.Error("Unable to save text data: failed to get user id from context",
			slog.String("caller", caller.CodeLine()))
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
		slog.Info("Unable to save text data: invalid request",
			slog.String("user_d", userID),
			slog.Any("violated_fields", report))
		return nil, processValidationError(report)
	}

	err := h.textDataService.Save(ctx, texData)
	if err != nil {
		slog.Error("Unable to save text data: internal error",
			slog.String("user_d", userID),
			slog.String("error", fmt.Errorf("%s %w", caller.CodeLine(), err).Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.PostTextDataResponse{}, nil
}
