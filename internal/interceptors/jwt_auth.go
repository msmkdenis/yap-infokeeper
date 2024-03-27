package interceptors

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-infokeeper/pkg/jwtgen"
)

var authMandatoryMethods = map[string]struct{}{
	"/proto.CreditCardService/PostSaveCreditCard": {},
}

type UserIDContextKey string

// JWTAuth represents JWT authentication interceptor.
type JWTAuth struct {
	jwtManager *jwtgen.JWTManager
}

func NewJWTAuth(jwtManager *jwtgen.JWTManager) *JWTAuth {
	return &JWTAuth{jwtManager: jwtManager}
}

// GRPCJWTAuth checks token from gRPC metadata and sets userID in the context. otherwise returns 401.
func (j *JWTAuth) GRPCJWTAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if _, ok := authMandatoryMethods[info.FullMethod]; !ok {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}

	c := md.Get(j.jwtManager.TokenName)
	if len(c) < 1 {
		slog.Info("authentification failed", slog.String("error", "no token found"))
		return nil, status.Errorf(codes.Unauthenticated, "no token found")
	}

	userID, err := j.jwtManager.GetUserID(c[0])
	if err != nil {
		slog.Info("authentification failed", slog.String("error", "authentification by UserID failed"))
		return nil, status.Errorf(codes.Unauthenticated, "authentification by UserID failed")
	}

	ctx = context.WithValue(ctx, UserIDContextKey("userID"), userID)
	return handler(ctx, req)
}
