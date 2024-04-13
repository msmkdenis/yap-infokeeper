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
	"github.com/msmkdenis/yap-infokeeper/internal/proto/credit_card"
)

func (c *CreditCardHandlerTestSuite) Test_GetLoadCreditCard() {
	token, err := c.jwtManager.BuildJWTString("050a289a-d10a-417b-ab89-3acfca0f6529")
	require.NoError(c.T(), err)

	cards := []model.CreditCard{
		{
			ID:        "050a289a-d10a-417b-ab89-3acfca0f6529",
			Number:    "1234 5678 9012 3456",
			OwnerID:   "050a289a-d10a-417b-ab89-3acfca0f6529",
			OwnerName: "John Doe",
			CVVCode:   "111",
			PinCode:   "5555",
			ExpiresAt: time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
			Metadata:  "some metadata",
			CreatedAt: time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:        "e6bee957-92ab-4733-82e7-df842ae27734",
			Number:    "4321 5678 9012 3456",
			OwnerID:   "050a289a-d10a-417b-ab89-3acfca0f6529",
			OwnerName: "John Doe",
			CVVCode:   "111",
			PinCode:   "5555",
			ExpiresAt: time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
			Metadata:  "some metadata",
			CreatedAt: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	testCases := []struct {
		name                  string
		token                 string
		body                  *credit_card.GetCreditCardRequest
		expectedCode          codes.Code
		expectedStatusMessage string
		prepare               func()
		expectedBody          []*credit_card.CreditCardCredentials
	}{
		{
			name:                  "Unauthorized - token not found",
			token:                 "",
			body:                  &credit_card.GetCreditCardRequest{},
			expectedCode:          codes.Unauthenticated,
			expectedStatusMessage: "authentification by UserID failed",
			prepare: func() {
				c.creditCardService.EXPECT().Load(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:                  "Internal - internal error",
			token:                 token,
			body:                  &credit_card.GetCreditCardRequest{},
			expectedCode:          codes.Internal,
			expectedStatusMessage: "internal error",
			prepare: func() {
				c.creditCardService.EXPECT().Load(gomock.Any(), gomock.Any()).Times(1).Return(nil, errors.New("error"))
			},
		},
		{
			name:  "Bad request - invalid date CreatedAfter",
			token: token,
			body: &credit_card.GetCreditCardRequest{
				CreatedAfter: "invalid date",
			},
			expectedCode:          codes.InvalidArgument,
			expectedStatusMessage: "date must be in format 2006-01-02",
			prepare: func() {
				c.creditCardService.EXPECT().Load(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Bad request - invalid date CreatedBefore",
			token: token,
			body: &credit_card.GetCreditCardRequest{
				CreatedBefore: "invalid date",
			},
			expectedCode:          codes.InvalidArgument,
			expectedStatusMessage: "date must be in format 2006-01-02",
			prepare: func() {
				c.creditCardService.EXPECT().Load(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Bad request - invalid date ExpiresAfter",
			token: token,
			body: &credit_card.GetCreditCardRequest{
				ExpiresAfter: "invalid date",
			},
			expectedCode:          codes.InvalidArgument,
			expectedStatusMessage: "date must be in format 2006-01-02",
			prepare: func() {
				c.creditCardService.EXPECT().Load(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Bad request - invalid date ExpireBefore",
			token: token,
			body: &credit_card.GetCreditCardRequest{
				ExpiresBefore: "invalid date",
			},
			expectedCode:          codes.InvalidArgument,
			expectedStatusMessage: "date must be in format 2006-01-02",
			prepare: func() {
				c.creditCardService.EXPECT().Load(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:         "Success - find all cards",
			token:        token,
			body:         &credit_card.GetCreditCardRequest{},
			expectedCode: codes.OK,
			prepare: func() {
				c.creditCardService.EXPECT().Load(gomock.Any(), gomock.Any()).Times(1).Return(cards, nil)
			},
			expectedBody: []*credit_card.CreditCardCredentials{
				{
					Number:    cards[0].Number,
					Owner:     cards[0].OwnerName,
					ExpiresAt: cards[0].ExpiresAt.Format("2006-01-02"),
					CvvCode:   cards[0].CVVCode,
					PinCode:   cards[0].PinCode,
					Metadata:  cards[0].Metadata,
					CreatedAt: cards[0].CreatedAt.Format("2006-01-02"),
				},
				{
					Number:    cards[1].Number,
					Owner:     cards[1].OwnerName,
					ExpiresAt: cards[1].ExpiresAt.Format("2006-01-02"),
					CvvCode:   cards[1].CVVCode,
					PinCode:   cards[1].PinCode,
					Metadata:  cards[1].Metadata,
					CreatedAt: cards[1].CreatedAt.Format("2006-01-02"),
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

			client := credit_card.NewCreditCardServiceClient(conn)
			resp, err := client.GetLoadCreditCard(ctx, test.body)
			if resp != nil {
				assert.Equal(t, test.expectedBody, resp.Cards)
			}

			st := status.Convert(err)
			assert.Equal(t, test.expectedCode, st.Code())
			assert.Equal(t, test.expectedStatusMessage, st.Message())
		})
	}
}
