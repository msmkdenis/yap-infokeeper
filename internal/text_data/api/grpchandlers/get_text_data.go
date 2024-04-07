package grpchandlers

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-infokeeper/internal/interceptors"
	pb "github.com/msmkdenis/yap-infokeeper/internal/text_data/api/grpchandlers/proto"
	"github.com/msmkdenis/yap-infokeeper/internal/text_data/specification"
)

func (h *TextData) GetLoadTextData(ctx context.Context, in *pb.GetTextDataRequest) (*pb.GetTextDataResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDContextKey("userID")).(string)
	if !ok {
		slog.Error("Internal server error", slog.String("error", "unable to get userID from context"))
		return nil, status.Error(codes.Internal, "internal error")
	}

	spec, err := specification.NewTextDataSpecification(userID, in)
	if err != nil {
		slog.Error("invalid text data request", slog.String("error", err.Error()))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	texData, err := h.textDataService.Load(ctx, spec)
	if err != nil {
		slog.Error("failed to load text data", slog.String("error", err.Error()))
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
