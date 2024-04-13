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

	"github.com/msmkdenis/yap-infokeeper/internal/proto/text_data"
)

func (c *TextDataHandlerTestSuite) Test_PostSaveTextData() {
	token, err := c.jwtManager.BuildJWTString("050a289a-d10a-417b-ab89-3acfca0f6529")
	require.NoError(c.T(), err)

	testCases := []struct {
		name                         string
		token                        string
		body                         *text_data.PostTextDataRequest
		expectedCode                 codes.Code
		expectedStatusMessage        string
		expectedViolationField       string
		expectedViolationDescription string
		prepare                      func()
	}{
		{
			name:  "BadRequest - invalid uuid",
			token: token,
			body: &text_data.PostTextDataRequest{
				Uuid:     "invalid uuid",
				Data:     "some data",
				Metadata: "some metadata",
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid text data request",
			expectedViolationField:       "ID",
			expectedViolationDescription: "must be valid uuid",
			prepare: func() {
				c.textDataService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "BadRequest - empty data",
			token: token,
			body: &text_data.PostTextDataRequest{
				Uuid:     "050a289a-d10a-417b-ab89-3acfca0f6529",
				Data:     "",
				Metadata: "some metadata",
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid text data request",
			expectedViolationField:       "Data",
			expectedViolationDescription: "must be not empty",
			prepare: func() {
				c.textDataService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Unauthorized - token not found",
			token: "",
			body: &text_data.PostTextDataRequest{
				Uuid:     "050a289a-d10a-417b-ab89-3acfca0f6529",
				Data:     "some data",
				Metadata: "some metadata",
			},
			expectedCode:          codes.Unauthenticated,
			expectedStatusMessage: "authentification by UserID failed",
			prepare: func() {
				c.textDataService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Internal error - unable to save text data",
			token: token,
			body: &text_data.PostTextDataRequest{
				Uuid:     "050a289a-d10a-417b-ab89-3acfca0f6529",
				Data:     "some data",
				Metadata: "some metadata",
			},
			expectedCode:          codes.Internal,
			expectedStatusMessage: "internal error",
			prepare: func() {
				c.textDataService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(1).Return(errors.New("some error"))
			},
		},
		{
			name:  "Successful text data saved",
			token: token,
			body: &text_data.PostTextDataRequest{
				Uuid:     "050a289a-d10a-417b-ab89-3acfca0f6529",
				Data:     "some data",
				Metadata: "some metadata",
			},
			expectedCode: codes.OK,
			prepare: func() {
				c.textDataService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(1).Return(nil)
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

			client := text_data.NewTextDataServiceClient(conn)
			_, err := client.PostSaveTextData(ctx, test.body)

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
