package grpchandlers

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"github.com/msmkdenis/yap-infokeeper/internal/config"
	"github.com/msmkdenis/yap-infokeeper/internal/interceptors"
	pbUser "github.com/msmkdenis/yap-infokeeper/internal/user/api/grpchandlers/proto"
	"github.com/msmkdenis/yap-infokeeper/internal/user/mocks"
	"github.com/msmkdenis/yap-infokeeper/internal/user/model"
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

func (u *UserHandlerTestSuite) Test_PostRegisterUser() {
	token, err := u.jwtManager.BuildJWTString("050a289a-d10a-417b-ab89-3acfca0f6529")
	require.NoError(u.T(), err)

	testCases := []struct {
		name                         string
		body                         *pbUser.PostUserRegisterRequest
		expectedCode                 codes.Code
		expectedStatusMessage        string
		expectedViolationField       string
		expectedViolationDescription string
		prepare                      func()
		expectedToken                string
		expectedResponse             *pbUser.PostUserRegisterResponse
	}{
		{
			name: "BadRequest - invalid uuid",
			body: &pbUser.PostUserRegisterRequest{
				Id:       "non-uuid",
				Login:    "msmkdenis@gmail.com",
				Password: []byte("test"),
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid user request",
			expectedViolationField:       "ID",
			expectedViolationDescription: "must be valid uuid",
			prepare: func() {
				u.userService.EXPECT().Register(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "BadRequest - invalid email",
			body: &pbUser.PostUserRegisterRequest{
				Id:       "050a289a-d10a-417b-ab89-3acfca0f6529",
				Login:    "invalid-email",
				Password: []byte("test"),
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid user request",
			expectedViolationField:       "Login",
			expectedViolationDescription: "must be valid email",
			prepare: func() {
				u.userService.EXPECT().Register(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "BadRequest - zero length password",
			body: &pbUser.PostUserRegisterRequest{
				Id:       "050a289a-d10a-417b-ab89-3acfca0f6529",
				Login:    "msmkdenis@gmail.com",
				Password: []byte(""),
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid user request",
			expectedViolationField:       "Password",
			expectedViolationDescription: "is required",
			prepare: func() {
				u.userService.EXPECT().Register(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "InternalError - failed to save user",
			body: &pbUser.PostUserRegisterRequest{
				Id:       "050a289a-d10a-417b-ab89-3acfca0f6529",
				Login:    "msmkdenis@gmail.com",
				Password: []byte("test"),
			},
			expectedCode:                 codes.Internal,
			expectedStatusMessage:        "internal error",
			expectedViolationField:       "",
			expectedViolationDescription: "",
			prepare: func() {
				u.userService.EXPECT().Register(gomock.Any(), gomock.Any()).Times(1).Return(fmt.Errorf("failed to save user"))
			},
		},
		{
			name: "Successful registration",
			body: &pbUser.PostUserRegisterRequest{
				Id:       "050a289a-d10a-417b-ab89-3acfca0f6529",
				Login:    "msmkdenis@gmail.com",
				Password: []byte("test"),
			},
			expectedCode:                 codes.OK,
			expectedStatusMessage:        "",
			expectedViolationField:       "",
			expectedViolationDescription: "",
			expectedToken:                token,
			prepare: func() {
				u.userService.EXPECT().Register(gomock.Any(), model.User{
					ID:       "050a289a-d10a-417b-ab89-3acfca0f6529",
					Login:    "msmkdenis@gmail.com",
					Password: []byte("test"),
				}).Times(1).Return(nil)
			},
		},
	}

	for _, test := range testCases {
		u.T().Run(test.name, func(t *testing.T) {
			test.prepare()

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
			if resp != nil {
				assert.Equal(t, test.expectedToken, resp.Token)
			}
			for _, detail := range st.Details() {
				switch d := detail.(type) { //nolint:gocritic
				case *errdetails.BadRequest:
					for _, violation := range d.GetFieldViolations() {
						assert.Equal(t, test.expectedViolationField, violation.GetField())
						assert.Equal(t, test.expectedViolationDescription, violation.GetDescription())
					}
				}
			}
		})
	}
}
