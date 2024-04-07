package grpchandlers

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/msmkdenis/yap-infokeeper/internal/credit_card/api/grpchandlers/proto"
	"github.com/msmkdenis/yap-infokeeper/internal/credit_card/specification"
	"github.com/msmkdenis/yap-infokeeper/internal/interceptors"
)

func (h *CreditCard) GetLoadCreditCard(ctx context.Context, in *pb.GetCreditCardRequest) (*pb.GetCreditCardResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDContextKey("userID")).(string)
	if !ok {
		slog.Error("Internal server error", slog.String("error", "unable to get userID from context"))
		return nil, status.Error(codes.Internal, "internal error")
	}

	spec, err := specification.NewCreditCardSpecification(userID, in)
	if err != nil {
		slog.Error("invalid text data request", slog.String("error", err.Error()))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	creditCards, err := h.creditCardService.Load(ctx, spec)
	if err != nil {
		slog.Error("failed to load credit cards", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	creditCardCredentials := make([]*pb.CreditCardCredentials, 0, len(creditCards))
	for _, v := range creditCards {
		creditCardCredentials = append(creditCardCredentials, &pb.CreditCardCredentials{
			Number:    v.Number,
			Owner:     v.OwnerName,
			ExpiresAt: v.ExpiresAt.Format("2006-01-02"),
			CvvCode:   v.CVVCode,
			PinCode:   v.PinCode,
			Metadata:  v.Metadata,
			CreatedAt: v.CreatedAt.Format("2006-01-02"),
		})
	}

	return &pb.GetCreditCardResponse{Cards: creditCardCredentials}, nil
}
