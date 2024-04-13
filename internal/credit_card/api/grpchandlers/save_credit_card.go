package grpchandlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-infokeeper/internal/interceptors"
	"github.com/msmkdenis/yap-infokeeper/internal/model"
	pb "github.com/msmkdenis/yap-infokeeper/internal/proto/credit_card"
	"github.com/msmkdenis/yap-infokeeper/pkg/caller"
)

func (h *CreditCard) PostSaveCreditCard(ctx context.Context, in *pb.PostCreditCardCredentialsRequest) (*pb.PostCreditCardCredentialsResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDContextKey("userID")).(string)
	if !ok {
		slog.Error("Unable to save credit card: failed to get user id from context",
			slog.String("caller", caller.CodeLine()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	expire, err := time.Parse("2006-01-02", in.ExpiresAt)
	if err != nil {
		slog.Error("Unable to save credit card: invalid credit card request",
			slog.String("user_d", userID),
			slog.String("caller", caller.CodeLine()))
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
	if errors.Is(err, model.ErrCardAlreadyExists) {
		slog.Info("Unable to to save credit card: credit card with number already exists",
			slog.String("user_id", userID),
			slog.String("card_number", card.Number),
			slog.String("error", fmt.Errorf("%s %w", caller.CodeLine(), err).Error()))
		return nil, status.Error(codes.AlreadyExists, "credit card already exists")
	}

	if err != nil {
		slog.Info("Unable to to save credit card: internal error",
			slog.String("user_id", userID),
			slog.String("error", fmt.Errorf("%s %w", caller.CodeLine(), err).Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.PostCreditCardCredentialsResponse{}, nil
}
