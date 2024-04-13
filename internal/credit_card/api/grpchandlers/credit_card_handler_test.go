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
	"github.com/msmkdenis/yap-infokeeper/internal/credit_card/mocks"
	"github.com/msmkdenis/yap-infokeeper/internal/interceptors"
	pbCreditCard "github.com/msmkdenis/yap-infokeeper/internal/proto/credit_card"
	"github.com/msmkdenis/yap-infokeeper/pkg/jwtgen"
)

var cfgMock = &config.Config{
	GRPCServer:    ":3300",
	TokenName:     "token",
	TokenSecret:   "secret",
	TokenExpHours: 24,
}

type CreditCardHandlerTestSuite struct {
	suite.Suite
	creditCardService *mocks.MockCreditCardService
	dialer            func(ctx context.Context, address string) (net.Conn, error)
	jwtManager        *jwtgen.JWTManager
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(CreditCardHandlerTestSuite))
}

func (c *CreditCardHandlerTestSuite) SetupSuite() {
	ctrl := gomock.NewController(c.T())
	c.creditCardService = mocks.NewMockCreditCardService(ctrl)
	c.jwtManager = jwtgen.NewJWTManager(cfgMock.TokenName, cfgMock.TokenSecret, cfgMock.TokenExpHours)
	auth := interceptors.NewJWTAuth(c.jwtManager).GRPCJWTAuth

	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(auth))
	pbCreditCard.RegisterCreditCardServiceServer(server, NewCreditCard(c.creditCardService))

	c.dialer = func(ctx context.Context, address string) (net.Conn, error) {
		return lis.Dial()
	}

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}
