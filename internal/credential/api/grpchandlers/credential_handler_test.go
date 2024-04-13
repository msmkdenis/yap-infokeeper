package grpchandlers

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/msmkdenis/yap-infokeeper/internal/config"
	"github.com/msmkdenis/yap-infokeeper/internal/credential/mocks"
	"github.com/msmkdenis/yap-infokeeper/internal/interceptors"
	pb "github.com/msmkdenis/yap-infokeeper/internal/proto/credential"
	"github.com/msmkdenis/yap-infokeeper/pkg/jwtgen"
)

var cfgMock = &config.Config{
	GRPCServer:    ":3300",
	TokenName:     "token",
	TokenSecret:   "secret",
	TokenExpHours: 24,
}

type CredentialHandlerTestSuite struct {
	suite.Suite
	credentialService *mocks.MockCredentialService
	dialer            func(ctx context.Context, address string) (net.Conn, error)
	jwtManager        *jwtgen.JWTManager
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(CredentialHandlerTestSuite))
}

func (c *CredentialHandlerTestSuite) SetupSuite() {
	ctrl := gomock.NewController(c.T())
	c.credentialService = mocks.NewMockCredentialService(ctrl)
	c.jwtManager = jwtgen.NewJWTManager(cfgMock.TokenName, cfgMock.TokenSecret, cfgMock.TokenExpHours)
	auth := interceptors.NewJWTAuth(c.jwtManager).GRPCJWTAuth

	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(auth))
	pb.RegisterCredentialServiceServer(server, NewCredential(c.credentialService))

	c.dialer = func(ctx context.Context, address string) (net.Conn, error) {
		return lis.Dial()
	}

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}
