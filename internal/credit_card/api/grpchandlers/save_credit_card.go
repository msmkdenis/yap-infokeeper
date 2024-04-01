package grpchandlers

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/msmkdenis/yap-infokeeper/internal/credit_card/api/grpchandlers/proto"
	"github.com/msmkdenis/yap-infokeeper/internal/interceptors"
	"github.com/msmkdenis/yap-infokeeper/internal/model"
	apperr "github.com/msmkdenis/yap-infokeeper/pkg/apperror"
)

func (h *CreditCard) PostSaveCreditCard(ctx context.Context, in *pb.PostCreditCardCredentialsRequest) (*pb.PostCreditCardCredentialsResponse, error) {
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
		OwnerID:   userID,
		OwnerName: in.Owner,
		ExpiresAt: expire,
		CVVCode:   in.CvvCode,
		PinCode:   in.PinCode,
		Metadata:  in.Metadata,
	}

	report, ok := h.validator.ValidateCreditCard(card)
	if !ok {
		return nil, processValidationError(report)
	}

	err = h.creditCardService.Save(ctx, card)
	if errors.Is(err, apperr.ErrCardAlreadyExists) {
		slog.Info("Credit card already exists", slog.String("with number", in.Number))
		return nil, status.Error(codes.InvalidArgument, "card with given number already exists")
	}

	if err != nil {
		slog.Info("Unable to save credit card", slog.String("with number", in.Number))
		return nil, status.Error(codes.Internal, "internal error while saving credit card")
	}

	return &pb.PostCreditCardCredentialsResponse{}, nil
}
