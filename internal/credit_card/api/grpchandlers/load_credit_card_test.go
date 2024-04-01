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

	pbCreditCard "github.com/msmkdenis/yap-infokeeper/internal/credit_card/api/grpchandlers/proto"
	"github.com/msmkdenis/yap-infokeeper/internal/model"
	apperr "github.com/msmkdenis/yap-infokeeper/pkg/apperror"
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
		},
	}

	testCases := []struct {
		name                  string
		token                 string
		body                  *pbCreditCard.GetCreditCardRequest
		expectedCode          codes.Code
		expectedStatusMessage string
		prepare               func()
		expectedBody          []*pbCreditCard.CreditCardCredentials
	}{
		{
			name:                  "Unauthorized - token not found",
			token:                 "",
			body:                  &pbCreditCard.GetCreditCardRequest{CardNumber: ""},
			expectedCode:          codes.Unauthenticated,
			expectedStatusMessage: "authentification by UserID failed",
			prepare: func() {
				c.creditCardService.EXPECT().SelectByOwnerIDCardNumber(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				c.creditCardService.EXPECT().SelectAllByOwnerID(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:                  "Internal - internal error find exact card",
			token:                 token,
			body:                  &pbCreditCard.GetCreditCardRequest{CardNumber: "1234 5678 9012 3456"},
			expectedCode:          codes.Internal,
			expectedStatusMessage: "internal error",
			prepare: func() {
				c.creditCardService.EXPECT().SelectByOwnerIDCardNumber(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil, errors.New("error"))
				c.creditCardService.EXPECT().SelectAllByOwnerID(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:                  "Internal - internal error find all cards",
			token:                 token,
			body:                  &pbCreditCard.GetCreditCardRequest{CardNumber: ""},
			expectedCode:          codes.Internal,
			expectedStatusMessage: "internal error",
			prepare: func() {
				c.creditCardService.EXPECT().SelectByOwnerIDCardNumber(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				c.creditCardService.EXPECT().SelectAllByOwnerID(gomock.Any(), gomock.Any()).Times(1).Return(nil, errors.New("error"))
			},
		},
		{
			name:         "Success - find all cards",
			token:        token,
			body:         &pbCreditCard.GetCreditCardRequest{CardNumber: ""},
			expectedCode: codes.OK,
			prepare: func() {
				c.creditCardService.EXPECT().SelectByOwnerIDCardNumber(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				c.creditCardService.EXPECT().SelectAllByOwnerID(gomock.Any(), gomock.Any()).Times(1).Return(cards, nil)
			},
			expectedBody: []*pbCreditCard.CreditCardCredentials{
				{
					Number:    cards[0].Number,
					Owner:     cards[0].OwnerName,
					ExpiresAt: cards[0].ExpiresAt.Format("2006-01-02"),
					CvvCode:   cards[0].CVVCode,
					PinCode:   cards[0].PinCode,
					Metadata:  cards[0].Metadata,
				},
				{
					Number:    cards[1].Number,
					Owner:     cards[1].OwnerName,
					ExpiresAt: cards[1].ExpiresAt.Format("2006-01-02"),
					CvvCode:   cards[1].CVVCode,
					PinCode:   cards[1].PinCode,
					Metadata:  cards[1].Metadata,
				},
			},
		},
		{
			name:         "Success - find exact card",
			token:        token,
			body:         &pbCreditCard.GetCreditCardRequest{CardNumber: "1234 5678 9012 3456"},
			expectedCode: codes.OK,
			prepare: func() {
				c.creditCardService.EXPECT().SelectByOwnerIDCardNumber(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&cards[0], nil)
				c.creditCardService.EXPECT().SelectAllByOwnerID(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedBody: []*pbCreditCard.CreditCardCredentials{
				{
					Number:    cards[0].Number,
					Owner:     cards[0].OwnerName,
					ExpiresAt: cards[0].ExpiresAt.Format("2006-01-02"),
					CvvCode:   cards[0].CVVCode,
					PinCode:   cards[0].PinCode,
					Metadata:  cards[0].Metadata,
				},
			},
		},
		{
			name:         "Success - return empty slice",
			token:        token,
			body:         &pbCreditCard.GetCreditCardRequest{CardNumber: "1234 5678 9012 3456"},
			expectedCode: codes.OK,
			prepare: func() {
				c.creditCardService.EXPECT().SelectByOwnerIDCardNumber(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil, apperr.ErrCardNotFound)
				c.creditCardService.EXPECT().SelectAllByOwnerID(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedBody: []*pbCreditCard.CreditCardCredentials(nil),
		},
		{
			name:         "Success - return empty slice",
			token:        token,
			body:         &pbCreditCard.GetCreditCardRequest{CardNumber: ""},
			expectedCode: codes.OK,
			prepare: func() {
				c.creditCardService.EXPECT().SelectByOwnerIDCardNumber(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				c.creditCardService.EXPECT().SelectAllByOwnerID(gomock.Any(), gomock.Any()).Times(1).Return(make([]model.CreditCard, 0), nil)
			},
			expectedBody: []*pbCreditCard.CreditCardCredentials(nil),
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

			client := pbCreditCard.NewCreditCardServiceClient(conn)
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
