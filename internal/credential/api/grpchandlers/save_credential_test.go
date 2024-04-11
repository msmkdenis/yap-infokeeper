package grpchandlers

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-infokeeper/internal/proto/credential"
)

func (c *CredentialHandlerTestSuite) Test_PostSaveCredential() {
	token, err := c.jwtManager.BuildJWTString("050a289a-d10a-417b-ab89-3acfca0f6529")
	require.NoError(c.T(), err)

	testCases := []struct {
		name                         string
		token                        string
		body                         *credential.PostCredentialRequest
		expectedCode                 codes.Code
		expectedStatusMessage        string
		expectedViolationField       string
		expectedViolationDescription string
		prepare                      func()
	}{
		{
			name:  "BadRequest - invalid uuid",
			token: token,
			body: &credential.PostCredentialRequest{
				Uuid:     "invalid uuid",
				Login:    "some login",
				Password: "qwerty",
				Metadata: "some data",
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid credential request",
			expectedViolationField:       "ID",
			expectedViolationDescription: "must be valid uuid",
			prepare: func() {
				c.credentialService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "BadRequest - empty login",
			token: token,
			body: &credential.PostCredentialRequest{
				Uuid:     "050a289a-d10a-417b-ab89-3acfca0f6529",
				Login:    "",
				Password: "qwerty",
				Metadata: "some data",
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid credential request",
			expectedViolationField:       "Login",
			expectedViolationDescription: "must be not empty",
			prepare: func() {
				c.credentialService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "BadRequest - empty password",
			token: token,
			body: &credential.PostCredentialRequest{
				Uuid:     "050a289a-d10a-417b-ab89-3acfca0f6529",
				Login:    "some login",
				Password: "",
				Metadata: "some data",
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid credential request",
			expectedViolationField:       "Password",
			expectedViolationDescription: "must be not empty",
			prepare: func() {
				c.credentialService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Unauthorized - token not found",
			token: "",
			body: &credential.PostCredentialRequest{
				Uuid:     "050a289a-d10a-417b-ab89-3acfca0f6529",
				Login:    "some login",
				Password: "qwerty",
				Metadata: "some data",
			},
			expectedCode:          codes.Unauthenticated,
			expectedStatusMessage: "authentification by UserID failed",
			prepare: func() {
				c.credentialService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Internal error - unable to save credential",
			token: token,
			body: &credential.PostCredentialRequest{
				Uuid:     "050a289a-d10a-417b-ab89-3acfca0f6529",
				Login:    "some login",
				Password: "qwerty",
				Metadata: "some data",
			},
			expectedCode:          codes.Internal,
			expectedStatusMessage: "internal error while saving credential",
			prepare: func() {
				c.credentialService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(1).Return(errors.New("internal"))
			},
		},
		{
			name:  "Successful credential saved",
			token: token,
			body: &credential.PostCredentialRequest{
				Uuid:     "050a289a-d10a-417b-ab89-3acfca0f6529",
				Login:    "some login",
				Password: "qwerty",
				Metadata: "some data",
			},
			expectedCode: codes.OK,
			prepare: func() {
				c.credentialService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
		},
	}
	for _, test := range testCases {
		c.T().Run(test.name, func(t *testing.T) {
			test.prepare()

			header := metadata.New(map[string]string{"token": test.token})
			ctx := metadata.NewOutgoingContext(context.Background(), header)
			conn, _ := grpc.DialContext(ctx, "bufnet",
				grpc.WithContextDialer(c.dialer),
				grpc.WithTransportCredentials(insecure.NewCredentials()))
			defer conn.Close()

			client := credential.NewCredentialServiceClient(conn)
			_, err := client.PostSaveCredential(ctx, test.body)

			st := status.Convert(err)
			assert.Equal(t, test.expectedCode, st.Code())
			assert.Equal(t, test.expectedStatusMessage, st.Message())
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
