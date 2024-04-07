package grpchandlers

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/msmkdenis/yap-infokeeper/internal/credential/api/grpchandlers/proto"
	"github.com/msmkdenis/yap-infokeeper/internal/credential/specification"
	"github.com/msmkdenis/yap-infokeeper/internal/interceptors"
)

func (h *Credential) GetLoadCredentials(ctx context.Context, in *pb.GetCredentialRequest) (*pb.GetCredentialResponse, error) {
	userID, ok := ctx.Value(interceptors.UserIDContextKey("userID")).(string)
	if !ok {
		slog.Error("Internal server error", slog.String("error", "unable to get userID from context"))
		return nil, status.Error(codes.Internal, "internal error")
	}

	credentialSpec, err := specification.NewCredentialSpecification(userID, in)
	if err != nil {
		slog.Error("invalid credential request", slog.String("error", err.Error()))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	creds, err := h.credentialService.Load(ctx, credentialSpec)
	if err != nil {
		slog.Error("failed to load credentials", slog.String("error", err.Error()))
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
