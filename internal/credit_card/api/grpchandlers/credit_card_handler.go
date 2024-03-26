package grpchandlers

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-infokeeper/internal/interceptors"

	pb "github.com/msmkdenis/yap-infokeeper/internal/credit_card/api/grpchandlers/proto"
	"github.com/msmkdenis/yap-infokeeper/internal/credit_card/model"
)

type CreditCardService interface {
	Save(ctx context.Context, ownerID string, card model.CreditCard) error
}

type UserRegister struct {
	creditCardService CreditCardService
	pb.UnimplementedCreditCardServiceServer
}

func NewCreditCardHandler(creditCardService CreditCardService) *UserRegister {
	return &UserRegister{creditCardService: creditCardService}
}

func (h *UserRegister) PostSaveCreditCard(ctx context.Context, in *pb.PostCreditCardCredentialsRequest) (*pb.PostCreditCardCredentialsResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDContextKey("userID")).(string)
	if !ok {
		slog.Error("Internal server error", slog.String("error", "unable to get userID from context"))
		return nil, status.Error(codes.Internal, "internal error")
	}

	expire, err := time.Parse("2006-01-02", in.ExpiresAt)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "date format must be 'YYYY-DD-MM'")
	}

	card := model.CreditCard{
		ID:        in.Uuid,
		Number:    in.Number,
		Owner:     userID,
		ExpiresAt: expire,
		CVVCode:   in.CvvCode,
		PinCode:   in.PinCode,
	}

	err = h.creditCardService.Save(ctx, userID, card)
	if err != nil {
		return nil, err
	}

	return &pb.PostCreditCardCredentialsResponse{}, nil
}
