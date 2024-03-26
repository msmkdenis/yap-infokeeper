package grpchandlers

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

func (u *UserHandlerTestSuite) SetupSuite() {
	ctrl := gomock.NewController(u.T())
	u.userService = mocks.NewMockUserService(ctrl)
	jwtManager := jwtgen.NewJWTManager(cfgMock.TokenName, cfgMock.TokenSecret, cfgMock.TokenExpHours)
	auth := interceptors.NewJWTAuth(jwtManager).GRPCJWTAuth

	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(auth))
	pbUser.RegisterUserServiceServer(server, NewUserRegister(u.userService, jwtManager))

	u.dialer = func(ctx context.Context, address string) (net.Conn, error) {
		return lis.Dial()
	}

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()
}

func (u *UserHandlerTestSuite) Test_PostRegisterUser() {
	testCases := []struct {
		name                  string
		body                  *pbUser.PostUserRegisterRequest
		expectedCode          codes.Code
		expectedBody          *pbUser.PostUserRegisterResponse
		expectedStatusMessage string
	}{
		{
			name:                  "BadRequest - invalid uuid",
			body:                  &pbUser.PostUserRegisterRequest{Id: "non-uuid", Login: "msmkdenis@gmail.com", Password: []byte("test")},
			expectedCode:          codes.InvalidArgument,
			expectedBody:          nil,
			expectedStatusMessage: "invalid user request",
		},
	}

	for _, test := range testCases {
		u.T().Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			conn, _ := grpc.DialContext(ctx, "bufnet",
				grpc.WithContextDialer(u.dialer),
				grpc.WithTransportCredentials(insecure.NewCredentials()))
			defer conn.Close()

			client := pbUser.NewUserServiceClient(conn)
			resp, err := client.PostRegisterUser(ctx, test.body)

			st := status.Convert(err)
			assert.Equal(t, test.expectedCode, st.Code())
			assert.Equal(t, test.expectedStatusMessage, st.Message())
			assert.Equal(t, test.expectedBody, resp)
		})
	}
}
