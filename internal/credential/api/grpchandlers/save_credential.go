package grpchandlers

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-infokeeper/internal/interceptors"
	"github.com/msmkdenis/yap-infokeeper/internal/model"
	pb "github.com/msmkdenis/yap-infokeeper/internal/proto/credential"
	"github.com/msmkdenis/yap-infokeeper/pkg/caller"
)

func (h *Credential) PostSaveCredential(ctx context.Context, in *pb.PostCredentialRequest) (*pb.PostCredentialResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDContextKey("userID")).(string)
	if !ok {
		slog.Error("Unable to save credential: failed to get user id from context",
			slog.String("caller", caller.CodeLine()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	credential := model.Credential{
		ID:       in.Uuid,
		OwnerID:  userID,
		Login:    in.Login,
		Password: in.Password,
		Metadata: in.Metadata,
	}

	report, ok := h.validator.ValidateCredential(credential)
	if !ok {
		return nil, processValidationError(report)
	}

	err := h.credentialService.Save(ctx, credential)
	if err != nil {
		slog.Info("Unable to to save credential: internal error",
			slog.String("user_id", userID),
			slog.String("error", fmt.Errorf("%s %w", caller.CodeLine(), err).Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.PostCredentialResponse{}, nil
}
