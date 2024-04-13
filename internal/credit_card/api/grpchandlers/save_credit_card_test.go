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

	"github.com/msmkdenis/yap-infokeeper/internal/proto/credit_card"
)

func (c *CreditCardHandlerTestSuite) Test_PostSaveCreditCard() {
	token, err := c.jwtManager.BuildJWTString("050a289a-d10a-417b-ab89-3acfca0f6529")
	require.NoError(c.T(), err)

	testCases := []struct {
		name                         string
		token                        string
		body                         *credit_card.PostCreditCardCredentialsRequest
		expectedCode                 codes.Code
		expectedStatusMessage        string
		expectedViolationField       string
		expectedViolationDescription string
		prepare                      func()
	}{
		{
			name:  "BadRequest - invalid uuid",
			token: token,
			body: &credit_card.PostCreditCardCredentialsRequest{
				Uuid:      "invalid uuid",
				Number:    "1234 5678 9012 3456",
				Owner:     "John Doe",
				ExpiresAt: "2025-12-01",
				CvvCode:   "111",
				PinCode:   "5555",
				Metadata:  "some metadata",
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid credit card request",
			expectedViolationField:       "ID",
			expectedViolationDescription: "must be valid uuid",
			prepare: func() {
				c.creditCardService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "BadRequest - invalid date",
			token: token,
			body: &credit_card.PostCreditCardCredentialsRequest{
				Uuid:      "050a289a-d10a-417b-ab89-3acfca0f6529",
				Number:    "1234 5678 9012 3456",
				Owner:     "John Doe",
				ExpiresAt: "2025-15-25",
				CvvCode:   "111",
				PinCode:   "5555",
				Metadata:  "some metadata",
			},
			expectedCode:          codes.InvalidArgument,
			expectedStatusMessage: "date format must be 'YYYY-DD-MM'",
			prepare: func() {
				c.creditCardService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "BadRequest - invalid number",
			token: token,
			body: &credit_card.PostCreditCardCredentialsRequest{
				Uuid:      "050a289a-d10a-417b-ab89-3acfca0f6529",
				Number:    "1234 5678 9012 3456 1234",
				Owner:     "John Doe",
				ExpiresAt: "2025-12-01",
				CvvCode:   "111",
				PinCode:   "5555",
				Metadata:  "some metadata",
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid credit card request",
			expectedViolationField:       "Number",
			expectedViolationDescription: "must be valid credit card number",
			prepare: func() {
				c.creditCardService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "BadRequest - invalid cvv",
			token: token,
			body: &credit_card.PostCreditCardCredentialsRequest{
				Uuid:      "050a289a-d10a-417b-ab89-3acfca0f6529",
				Number:    "1234 5678 9012 3456",
				Owner:     "John Doe",
				ExpiresAt: "2025-12-01",
				CvvCode:   "ABC",
				PinCode:   "5555",
				Metadata:  "some metadata",
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid credit card request",
			expectedViolationField:       "CVVCode",
			expectedViolationDescription: "must be valid cvv",
			prepare: func() {
				c.creditCardService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "BadRequest - invalid pin",
			token: token,
			body: &credit_card.PostCreditCardCredentialsRequest{
				Uuid:      "050a289a-d10a-417b-ab89-3acfca0f6529",
				Number:    "1234 5678 9012 3456",
				Owner:     "John Doe",
				ExpiresAt: "2025-12-01",
				CvvCode:   "123",
				PinCode:   "55AB",
				Metadata:  "some metadata",
			},
			expectedCode:                 codes.InvalidArgument,
			expectedStatusMessage:        "invalid credit card request",
			expectedViolationField:       "PinCode",
			expectedViolationDescription: "must be valid pin",
			prepare: func() {
				c.creditCardService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Unauthorized - token not found",
			token: "",
			body: &credit_card.PostCreditCardCredentialsRequest{
				Uuid:      "050a289a-d10a-417b-ab89-3acfca0f6529",
				Number:    "1234 5678 9012 3456",
				Owner:     "John Doe",
				ExpiresAt: "2025-12-01",
				CvvCode:   "123",
				PinCode:   "5555",
				Metadata:  "some metadata",
			},
			expectedCode:          codes.Unauthenticated,
			expectedStatusMessage: "authentification by UserID failed",
			prepare: func() {
				c.creditCardService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Internal error - unable to save credit card",
			token: token,
			body: &credit_card.PostCreditCardCredentialsRequest{
				Uuid:      "050a289a-d10a-417b-ab89-3acfca0f6529",
				Number:    "1234 5678 9012 3456",
				Owner:     "John Doe",
				ExpiresAt: "2025-12-01",
				CvvCode:   "123",
				PinCode:   "5555",
				Metadata:  "some metadata",
			},
			expectedCode:          codes.Internal,
			expectedStatusMessage: "internal error",
			prepare: func() {
				c.creditCardService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(1).Return(errors.New("internal"))
			},
		},
		{
			name:  "Successful credit card saved",
			token: token,
			body: &credit_card.PostCreditCardCredentialsRequest{
				Uuid:      "050a289a-d10a-417b-ab89-3acfca0f6529",
				Number:    "1234 5678 9012 3456",
				Owner:     "John Doe",
				ExpiresAt: "2025-12-01",
				CvvCode:   "123",
				PinCode:   "5555",
				Metadata:  "some metadata",
			},
			expectedCode: codes.OK,
			prepare: func() {
				c.creditCardService.EXPECT().Save(gomock.Any(), gomock.Any()).Times(1).Return(nil)
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

			client := credit_card.NewCreditCardServiceClient(conn)
			_, err := client.PostSaveCreditCard(ctx, test.body)

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
