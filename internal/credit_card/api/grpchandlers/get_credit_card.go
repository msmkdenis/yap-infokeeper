package grpchandlers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/msmkdenis/yap-infokeeper/pkg/caller"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-infokeeper/internal/credit_card/specification"
	"github.com/msmkdenis/yap-infokeeper/internal/interceptors"
	pb "github.com/msmkdenis/yap-infokeeper/internal/proto/credit_card"
)

func (h *CreditCard) GetLoadCreditCard(ctx context.Context, in *pb.GetCreditCardRequest) (*pb.GetCreditCardResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDContextKey("userID")).(string)
	if !ok {
		slog.Error("Unable to load credit card: failed to get user id from context",
			slog.String("caller", caller.CodeLine()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	spec, err := specification.NewCreditCardSpecification(userID, in)
	if err != nil {
		slog.Error("Unable to load credit card: invalid credit card request",
			slog.String("user_d", userID),
			slog.String("error", fmt.Errorf("%s %w", caller.CodeLine(), err).Error()))
		return nil, status.Error(codes.InvalidArgument, "date must be in format 2006-01-02")
	}

	creditCards, err := h.creditCardService.Load(ctx, spec)
	if err != nil {
		slog.Error("Unable to load credit card: internal error",
			slog.String("user_d", userID),
			slog.String("error", fmt.Errorf("%s %w", caller.CodeLine(), err).Error()))
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
