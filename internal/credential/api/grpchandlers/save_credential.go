package grpchandlers

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/msmkdenis/yap-infokeeper/internal/credential/api/grpchandlers/proto"
	"github.com/msmkdenis/yap-infokeeper/internal/interceptors"
	"github.com/msmkdenis/yap-infokeeper/internal/model"
)

func (h *Credential) PostSaveCredential(ctx context.Context, in *pb.PostCredentialRequest) (*pb.PostCredentialResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDContextKey("userID")).(string)
	if !ok {
		slog.Error("Internal server error", slog.String("error", "unable to get userID from context"))
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
		slog.Info("Unable to save credential", slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, "internal error while saving credential")
	}

	return &pb.PostCredentialResponse{}, nil
}
