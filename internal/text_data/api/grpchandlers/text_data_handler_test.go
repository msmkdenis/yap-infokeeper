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
	pb "github.com/msmkdenis/yap-infokeeper/internal/text_data/api/grpchandlers/proto"
	"github.com/msmkdenis/yap-infokeeper/internal/text_data/mocks"
	"github.com/msmkdenis/yap-infokeeper/pkg/jwtgen"
)

var cfgMock = &config.Config{
	GRPCServer:    ":3300",
	TokenName:     "token",
	TokenSecret:   "secret",
	TokenExpHours: 24,
}

type TextDataHandlerTestSuite struct {
	suite.Suite
	textDataService *mocks.MockTextDataService
	dialer          func(ctx context.Context, address string) (net.Conn, error)
	jwtManager      *jwtgen.JWTManager
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TextDataHandlerTestSuite))
}

func (c *TextDataHandlerTestSuite) SetupSuite() {
	ctrl := gomock.NewController(c.T())
	c.textDataService = mocks.NewMockTextDataService(ctrl)
	c.jwtManager = jwtgen.NewJWTManager(cfgMock.TokenName, cfgMock.TokenSecret, cfgMock.TokenExpHours)
	auth := interceptors.NewJWTAuth(c.jwtManager).GRPCJWTAuth

	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(auth))
	pb.RegisterTextDataServiceServer(server, NewTextData(c.textDataService))

	c.dialer = func(ctx context.Context, address string) (net.Conn, error) {
		return lis.Dial()
	}

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}
