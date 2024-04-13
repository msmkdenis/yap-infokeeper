package grpchandlers

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-infokeeper/internal/credential/specification"
	"github.com/msmkdenis/yap-infokeeper/internal/interceptors"
	pb "github.com/msmkdenis/yap-infokeeper/internal/proto/credential"
	"github.com/msmkdenis/yap-infokeeper/pkg/caller"
)

func (h *Credential) GetLoadCredentials(ctx context.Context, in *pb.GetCredentialRequest) (*pb.GetCredentialResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDContextKey("userID")).(string)
	if !ok {
		slog.Error("Unable to load credential: failed to get user id from context",
			slog.String("caller", caller.CodeLine()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	credentialSpec, err := specification.NewCredentialSpecification(userID, in)
	if err != nil {
		slog.Error("Unable to load credential: invalid credential request",
			slog.String("user_d", userID),
			slog.String("error", fmt.Errorf("%s %w", caller.CodeLine(), err).Error()))
		return nil, status.Error(codes.InvalidArgument, "date must be in format 2006-01-02")
	}

	creds, err := h.credentialService.Load(ctx, credentialSpec)
	if err != nil {
		slog.Error("Unable to load credential: internal error",
			slog.String("user_d", userID),
			slog.String("error", fmt.Errorf("%s %w", caller.CodeLine(), err).Error()))
		return nil, status.Error(codes.Internal, "internal error")
	}

	credentials := make([]*pb.Credential, 0, len(creds))
	for _, v := range creds {
		credentials = append(credentials, &pb.Credential{
			Login:     v.Login,
			Password:  v.Password,
			CreatedAt: v.CreatedAt.Format("2006-01-02"),
			Metadata:  v.Metadata,
		})
	}

	return &pb.GetCredentialResponse{Cards: credentials}, nil
}
