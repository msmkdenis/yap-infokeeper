package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
)

func (s *TextDataRepositoryTestSuite) Test_Insert() {
	textDataID, err := uuid.NewUUID()
	require.NoError(s.T(), err)

	ownerID, err := uuid.NewUUID()
	require.NoError(s.T(), err)

	textData := model.TextData{
		ID:       textDataID.String(),
		OwnerID:  ownerID.String(),
		Data:     "test data",
		Metadata: "test meta data",
	}

	testCases := []struct {
		name                  string
		textData              model.TextData
		expectedSavedTextData model.TextData
		fixtureFunc           func()
		expectedError         error
	}{
		{
			name:     "success",
			textData: textData,
			fixtureFunc: func() {
				s.insertTestUser(ownerID.String(), "login test", "password test")
			},
			expectedSavedTextData: textData,
			expectedError:         nil,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.fixtureFunc != nil {
				tc.fixtureFunc()
			}

			err := s.textDataRepository.Insert(context.Background(), tc.textData)
			require.Equal(s.T(), tc.expectedError, err)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			query := "select id, owner_id, data, created_at, metadata from infokeeper.text_data where id = $1"
			queryRows, err := s.pool.DB.Query(ctx, query, tc.textData.ID)
			assert.NoError(s.T(), err)

			savedData, err := pgx.CollectRows(queryRows, pgx.RowToStructByPos[model.TextData])
			assert.NoError(s.T(), err)

			assert.Equal(s.T(), tc.textData.ID, savedData[0].ID)
			assert.Equal(s.T(), tc.textData.OwnerID, savedData[0].OwnerID)
			assert.Equal(s.T(), tc.textData.Data, savedData[0].Data)
			assert.Equal(s.T(), tc.textData.Metadata, savedData[0].Metadata)
		})
	}
}

func (s *TextDataRepositoryTestSuite) insertTestUser(ownerID, login, password string) {
	query := "insert into infokeeper.user (id, login, password) values ($1, $2, $3)"
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(s.T(), err)

	_, err = s.pool.DB.Exec(context.Background(), query, ownerID, login, pass)
	assert.NoError(s.T(), err)
}
