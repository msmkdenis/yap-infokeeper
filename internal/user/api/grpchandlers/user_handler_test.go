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
	"github.com/msmkdenis/yap-infokeeper/internal/interceptors"
	pbUser "github.com/msmkdenis/yap-infokeeper/internal/user/api/grpchandlers/proto"
	"github.com/msmkdenis/yap-infokeeper/internal/user/mocks"
	"github.com/msmkdenis/yap-infokeeper/pkg/jwtgen"
)

var cfgMock = &config.Config{
	GRPCServer:    ":3300",
	TokenName:     "token",
	TokenSecret:   "secret",
	TokenExpHours: 24,
}

type UserHandlerTestSuite struct {
	suite.Suite
	userService *mocks.MockUserService
	dialer      func(ctx context.Context, address string) (net.Conn, error)
	jwtManager  *jwtgen.JWTManager
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

func (u *UserHandlerTestSuite) SetupSuite() {
	ctrl := gomock.NewController(u.T())
	u.userService = mocks.NewMockUserService(ctrl)
	u.jwtManager = jwtgen.NewJWTManager(cfgMock.TokenName, cfgMock.TokenSecret, cfgMock.TokenExpHours)
	auth := interceptors.NewJWTAuth(u.jwtManager).GRPCJWTAuth

	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(auth))
	pbUser.RegisterUserServiceServer(server, NewUserRegister(u.userService, u.jwtManager))

	u.dialer = func(ctx context.Context, address string) (net.Conn, error) {
		return lis.Dial()
	}

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}
