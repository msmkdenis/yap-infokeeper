package grpchandlers

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
	"github.com/msmkdenis/yap-infokeeper/internal/proto/credential"
)

func (c *CredentialHandlerTestSuite) Test_GetLoadCredentials() {
	token, err := c.jwtManager.BuildJWTString("050a289a-d10a-417b-ab89-3acfca0f6529")
	require.NoError(c.T(), err)

	creds := []model.Credential{
		{
			ID:        "050a289a-d10a-417b-ab89-3acfca0f6529",
			OwnerID:   "050a289a-d10a-417b-ab89-3acfca0f6529",
			Login:     "some login",
			Password:  "qwerty",
			Metadata:  "some data",
			CreatedAt: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:        "050a289a-d10a-417b-ab89-3acfca0f6525",
			OwnerID:   "050a289a-d10a-417b-ab89-3acfca0f6529",
			Login:     "another login",
			Password:  "ytrewq",
			Metadata:  "some data",
			CreatedAt: time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	testCases := []struct {
		name                  string
		token                 string
		body                  *credential.GetCredentialRequest
		expectedCode          codes.Code
		expectedStatusMessage string
		prepare               func()
		expectedBody          []*credential.Credential
	}{
		{
			name:                  "Unauthorized - token not found",
			token:                 "",
			body:                  &credential.GetCredentialRequest{},
			expectedCode:          codes.Unauthenticated,
			expectedStatusMessage: "authentification by UserID failed",
			prepare: func() {
				c.credentialService.EXPECT().Load(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Internal - internal error",
			token: token,
			body: &credential.GetCredentialRequest{
				Login:         "some login",
				Password:      "qwerty",
				Metadata:      "some data",
				CreatedAfter:  "2020-01-01",
				CreatedBefore: "2026-01-01",
			},
			expectedCode:          codes.Internal,
			expectedStatusMessage: "internal error",
			prepare: func() {
				c.credentialService.EXPECT().Load(gomock.Any(), gomock.Any()).Times(1).Return(nil, errors.New("error"))
			},
		},
		{
			name:  "Bad request - invalid date CreatedAfter",
			token: token,
			body: &credential.GetCredentialRequest{
				Login:         "some login",
				Password:      "qwerty",
				Metadata:      "some data",
				CreatedAfter:  "invalid date",
				CreatedBefore: "2026-01-01",
			},
			expectedCode:          codes.InvalidArgument,
			expectedStatusMessage: "date must be in format 2006-01-02",
			prepare: func() {
				c.credentialService.EXPECT().Load(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Bad request - invalid date CreatedBefore",
			token: token,
			body: &credential.GetCredentialRequest{
				Login:         "some login",
				Password:      "qwerty",
				Metadata:      "some data",
				CreatedAfter:  "2020-01-01",
				CreatedBefore: "invalid date",
			},
			expectedCode:          codes.InvalidArgument,
			expectedStatusMessage: "date must be in format 2006-01-02",
			prepare: func() {
				c.credentialService.EXPECT().Load(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Success - find credentials",
			token: token,
			body: &credential.GetCredentialRequest{
				Login:         "some login",
				Password:      "qwerty",
				Metadata:      "some data",
				CreatedAfter:  "2020-01-01",
				CreatedBefore: "2026-01-01",
			},
			expectedCode: codes.OK,
			prepare: func() {
				c.credentialService.EXPECT().Load(gomock.Any(), gomock.Any()).Times(1).Return(creds, nil)
			},
			expectedBody: []*credential.Credential{
				{
					Login:     creds[0].Login,
					Password:  creds[0].Password,
					Metadata:  creds[0].Metadata,
					CreatedAt: creds[0].CreatedAt.Format("2006-01-02"),
				},
				{
					Login:     creds[1].Login,
					Password:  creds[1].Password,
					Metadata:  creds[1].Metadata,
					CreatedAt: creds[1].CreatedAt.Format("2006-01-02"),
				},
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
			resp, err := client.GetLoadCredentials(ctx, test.body)
			if resp != nil {
				assert.Equal(t, test.expectedBody, resp.Cards)
			}

			st := status.Convert(err)
			assert.Equal(t, test.expectedCode, st.Code())
			assert.Equal(t, test.expectedStatusMessage, st.Message())
		})
	}
}
