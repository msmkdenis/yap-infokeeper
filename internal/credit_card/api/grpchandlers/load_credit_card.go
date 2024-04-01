package grpchandlers

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/msmkdenis/yap-infokeeper/internal/credit_card/api/grpchandlers/proto"
	"github.com/msmkdenis/yap-infokeeper/internal/interceptors"
)

func (h *CreditCard) GetLoadCreditCard(ctx context.Context, in *pb.GetCreditCardRequest) (*pb.GetCreditCardResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDContextKey("userID")).(string)
	if !ok {
		slog.Error("Internal server error", slog.String("error", "unable to get userID from context"))
		return nil, status.Error(codes.Internal, "internal error")
	}

	var creditCardCredentials []*pb.CreditCardCredentials
	if len(in.CardNumber) == 0 {
		cards, err := h.creditCardService.SelectAllByOwnerID(ctx, userID)
		if err != nil {
			slog.Error("Internal server error", slog.String("error", err.Error()))
			return nil, status.Error(codes.Internal, "internal error")
		}
		for _, v := range cards {
			creditCardCredentials = append(creditCardCredentials, &pb.CreditCardCredentials{
				Number:    v.Number,
				Owner:     v.OwnerName,
				ExpiresAt: v.ExpiresAt.Format("2006-01-02"),
				CvvCode:   v.CVVCode,
				PinCode:   v.PinCode,
				Metadata:  v.Metadata,
			})
		}
	} else {
		card, err := h.creditCardService.SelectByOwnerIDCardNumber(ctx, userID, in.CardNumber)
		if err != nil {
			slog.Error("Internal server error", slog.String("error", err.Error()))
			return nil, status.Error(codes.Internal, "internal error")
		}
		creditCardCredentials = append(creditCardCredentials, &pb.CreditCardCredentials{
			Number:    card.Number,
			Owner:     card.OwnerName,
			ExpiresAt: card.ExpiresAt.Format("2006-01-02"),
			CvvCode:   card.CVVCode,
			PinCode:   card.PinCode,
			Metadata:  card.Metadata,
		})
	}

	return &pb.GetCreditCardResponse{Cards: creditCardCredentials}, nil
}
