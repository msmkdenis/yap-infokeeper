package grpchandlers

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-infokeeper/internal/interceptors"
	pb "github.com/msmkdenis/yap-infokeeper/internal/proto/text_data"
	"github.com/msmkdenis/yap-infokeeper/internal/text_data/specification"
	"github.com/msmkdenis/yap-infokeeper/pkg/caller"
)

func (h *TextData) GetLoadTextData(ctx context.Context, in *pb.GetTextDataRequest) (*pb.GetTextDataResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDContextKey("userID")).(string)
	if !ok {
		slog.Error("Unable to load text data: failed to get user id from context",
			slog.String("caller", caller.CodeLine()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	spec, err := specification.NewTextDataSpecification(userID, in)
	if err != nil {
		slog.Error("Unable to load text data: invalid text data request",
			slog.String("user_d", userID),
			slog.String("error", fmt.Errorf("%s %w", caller.CodeLine(), err).Error()))
		return nil, status.Error(codes.InvalidArgument, "date must be in format '2006-01-02'")
	}

	texData, err := h.textDataService.Load(ctx, spec)
	if err != nil {
		slog.Error("Unable to load text data: internal error",
			slog.String("user_d", userID),
			slog.String("error", fmt.Errorf("%s %w", caller.CodeLine(), err).Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	textData := make([]*pb.TextData, 0, len(texData))
	for _, v := range texData {
		textData = append(textData, &pb.TextData{
			Data:      v.Data,
			Metadata:  v.Metadata,
			CreatedAt: v.CreatedAt.Format("2006-01-02"),
		})
	}

	return &pb.GetTextDataResponse{Data: textData}, nil
}
