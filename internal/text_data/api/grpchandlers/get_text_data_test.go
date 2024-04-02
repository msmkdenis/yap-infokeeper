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
	pb "github.com/msmkdenis/yap-infokeeper/internal/text_data/api/grpchandlers/proto"
)

func (c *TextDataHandlerTestSuite) Test_GetLoadCredentials() {
	token, err := c.jwtManager.BuildJWTString("050a289a-d10a-417b-ab89-3acfca0f6529")
	require.NoError(c.T(), err)

	textDataChunks := []model.TextData{
		{
			ID:        "050a289a-d10a-417b-ab89-3acfca0f6529",
			OwnerID:   "050a289a-d10a-417b-ab89-3acfca0f6529",
			Data:      "some data",
			Metadata:  "some metadata",
			CreatedAt: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:        "050a289a-d10a-417b-ab89-3acfca0f5256",
			OwnerID:   "050a289a-d10a-417b-ab89-3acfca0f6529",
			Data:      "another data",
			Metadata:  "another metadata",
			CreatedAt: time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	testCases := []struct {
		name                  string
		token                 string
		body                  *pb.GetTextDataRequest
		expectedCode          codes.Code
		expectedStatusMessage string
		prepare               func()
		expectedBody          []*pb.TextData
	}{
		{
			name:                  "Unauthorized - token not found",
			token:                 "",
			body:                  &pb.GetTextDataRequest{},
			expectedCode:          codes.Unauthenticated,
			expectedStatusMessage: "authentification by UserID failed",
			prepare: func() {
				c.textDataService.EXPECT().Load(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Internal - internal error",
			token: token,
			body: &pb.GetTextDataRequest{
				Data:          "data",
				Metadata:      "metadata",
				CreatedAfter:  "2020-01-01",
				CreatedBefore: "2026-01-01",
			},
			expectedCode:          codes.Internal,
			expectedStatusMessage: "internal error",
			prepare: func() {
				c.textDataService.EXPECT().Load(gomock.Any(), gomock.Any()).Times(1).Return(nil, errors.New("some error"))
			},
		},
		{
			name:  "Bad request - invalid date CreatedAfter",
			token: token,
			body: &pb.GetTextDataRequest{
				Data:          "data",
				Metadata:      "metadata",
				CreatedAfter:  "invalid date",
				CreatedBefore: "2026-01-01",
			},
			expectedCode:          codes.InvalidArgument,
			expectedStatusMessage: "created after must be in format '2006-01-02'",
			prepare: func() {
				c.textDataService.EXPECT().Load(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Bad request - invalid date CreatedBefore",
			token: token,
			body: &pb.GetTextDataRequest{
				Data:          "data",
				Metadata:      "metadata",
				CreatedAfter:  "2020-01-01",
				CreatedBefore: "invalid date",
			},
			expectedCode:          codes.InvalidArgument,
			expectedStatusMessage: "created before must be in format '2006-01-02'",
			prepare: func() {
				c.textDataService.EXPECT().Load(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			name:  "Success - find text data",
			token: token,
			body: &pb.GetTextDataRequest{
				Data:          "data",
				Metadata:      "metadata",
				CreatedAfter:  "2020-01-01",
				CreatedBefore: "2026-01-01",
			},
			expectedCode: codes.OK,
			prepare: func() {
				c.textDataService.EXPECT().Load(gomock.Any(), gomock.Any()).Times(1).Return(textDataChunks, nil)
			},
			expectedBody: []*pb.TextData{
				{
					Data:      textDataChunks[0].Data,
					Metadata:  textDataChunks[0].Metadata,
					CreatedAt: textDataChunks[0].CreatedAt.Format("2006-01-02"),
				},
				{
					Data:      textDataChunks[1].Data,
					Metadata:  textDataChunks[1].Metadata,
					CreatedAt: textDataChunks[1].CreatedAt.Format("2006-01-02"),
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

			client := pb.NewTextDataServiceClient(conn)
			resp, err := client.GetLoadTextData(ctx, test.body)
			if resp != nil {
				assert.Equal(t, test.expectedBody, resp.Data)
			}

			st := status.Convert(err)
			assert.Equal(t, test.expectedCode, st.Code())
			assert.Equal(t, test.expectedStatusMessage, st.Message())
		})
	}
}
