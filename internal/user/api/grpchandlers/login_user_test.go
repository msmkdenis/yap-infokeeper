package grpchandlers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
	"github.com/msmkdenis/yap-infokeeper/internal/proto/user"
	apperr "github.com/msmkdenis/yap-infokeeper/pkg/apperror"
)

func (u *UserHandlerTestSuite) Test_PostLoginUser() {
	token, err := u.jwtManager.BuildJWTString("050a289a-d10a-417b-ab89-3acfca0f6529")
	require.NoError(u.T(), err)

	testCases := []struct {
		name                         string
		body                         *user.PostUserLoginRequest
		expectedCode                 codes.Code
		expectedStatusMessage        string
		expectedViolationField       string
		expectedViolationDescription string
		prepare                      func()
		expectedToken                string
		expectedResponse             *user.PostUserLoginResponse
	}{
		{
			name: "BadRequest - invalid email",
			body: &user.PostUserLoginRequest{
				Login:    "invalid-email",
				Password: []byte("test"),
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid user request",
			expectedViolationField:       "Login",
			expectedViolationDescription: "must be valid email",
			prepare: func() {
				u.userService.EXPECT().Login(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "BadRequest - zero length password",
			body: &user.PostUserLoginRequest{
				Login:    "msmkdenis@gmail.com",
				Password: []byte(""),
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid user request",
			expectedViolationField:       "Password",
			expectedViolationDescription: "is required",
			prepare: func() {
				u.userService.EXPECT().Login(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name: "Unauthenticated - failed to login, invalid password",
			body: &user.PostUserLoginRequest{
				Login:    "msmkdenis@gmail.com",
				Password: []byte("test"),
			},
			expectedCode:                 codes.Unauthenticated,
			expectedStatusMessage:        "invalid password",
			expectedViolationField:       "",
			expectedViolationDescription: "",
			prepare: func() {
				u.userService.EXPECT().Login(gomock.Any(), gomock.Any()).Times(1).Return(nil, apperr.ErrInvalidPassword)
			},
		},
		{
			name: "Unauthenticated - failed to login, user not found",
			body: &user.PostUserLoginRequest{
				Login:    "msmkdenis@gmail.com",
				Password: []byte("test"),
			},
			expectedCode:                 codes.Unauthenticated,
			expectedStatusMessage:        "user not found",
			expectedViolationField:       "",
			expectedViolationDescription: "",
			prepare: func() {
				u.userService.EXPECT().Login(gomock.Any(), gomock.Any()).Times(1).Return(nil, apperr.ErrUserNotFound)
			},
		},
		{
			name: "Unauthenticated - internal error",
			body: &user.PostUserLoginRequest{
				Login:    "msmkdenis@gmail.com",
				Password: []byte("test"),
			},
			expectedCode:                 codes.Internal,
			expectedStatusMessage:        "internal error",
			expectedViolationField:       "",
			expectedViolationDescription: "",
			prepare: func() {
				u.userService.EXPECT().Login(gomock.Any(), gomock.Any()).Times(1).Return(nil, errors.New("internal error"))
			},
		},
		{
			name: "Successful login",
			body: &user.PostUserLoginRequest{
				Login:    "msmkdenis@gmail.com",
				Password: []byte("test"),
			},
			expectedCode:                 codes.OK,
			expectedStatusMessage:        "",
			expectedViolationField:       "",
			expectedViolationDescription: "",
			expectedToken:                token,
			prepare: func() {
				u.userService.EXPECT().Login(gomock.Any(), gomock.Any()).Times(1).Return(&model.User{
					ID:        "050a289a-d10a-417b-ab89-3acfca0f6529",
					Login:     "msmkdenis@gmail.com",
					Password:  []byte("test"),
					CreatedAt: time.Now(),
				}, nil)
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

			client := user.NewUserServiceClient(conn)
			resp, err := client.PostLoginUser(ctx, test.body)

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
