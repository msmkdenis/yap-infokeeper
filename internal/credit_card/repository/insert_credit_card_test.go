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

func (s *CreditCardRepositoryTestSuite) Test_Insert() {
	ownerID, err := uuid.NewUUID()
	require.NoError(s.T(), err)

	creditCardID, err := uuid.NewUUID()
	require.NoError(s.T(), err)

	creditCard := model.CreditCard{
		ID:        creditCardID.String(),
		Number:    "1234 5678 9012 3456",
		OwnerName: "John Doe",
		OwnerID:   ownerID.String(),
		ExpiresAt: time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
		CVVCode:   "111",
		PinCode:   "5555",
		Metadata:  "some metadata",
	}
	testCases := []struct {
		name                    string
		creditCard              model.CreditCard
		expectedSavedCreditCard model.CreditCard
		fixtureFunc             func()
		expectedError           error
	}{
		{
			name:       "Success",
			creditCard: creditCard,
			fixtureFunc: func() {
				s.insertTestUser(ownerID.String(), "login test", "password test")
			},
			expectedSavedCreditCard: creditCard,
			expectedError:           nil,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.fixtureFunc != nil {
				tc.fixtureFunc()
			}

			err := s.creditCardRepository.Insert(context.Background(), tc.creditCard)
			require.Equal(s.T(), tc.expectedError, err)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			query := `select
					id,
					number,
					owner_id,
					owner_name,
					expires_at,
					cvv_code,
					pin_code,
					created_at,
					metadata
				from infokeeper.credit_card
				where owner_id = $1`

			queryRows, err := s.pool.DB.Query(ctx, query, tc.creditCard.OwnerID)
			assert.NoError(s.T(), err)

			savedData, err := pgx.CollectRows(queryRows, pgx.RowToStructByPos[model.CreditCard])
			assert.NoError(s.T(), err)

			assert.Equal(s.T(), tc.expectedSavedCreditCard.ID, savedData[0].ID)
			assert.Equal(s.T(), tc.expectedSavedCreditCard.Number, savedData[0].Number)
			assert.Equal(s.T(), tc.expectedSavedCreditCard.OwnerID, savedData[0].OwnerID)
			assert.Equal(s.T(), tc.expectedSavedCreditCard.OwnerName, savedData[0].OwnerName)
			assert.Equal(s.T(), tc.expectedSavedCreditCard.ExpiresAt, savedData[0].ExpiresAt)
			assert.Equal(s.T(), tc.expectedSavedCreditCard.CVVCode, savedData[0].CVVCode)
			assert.Equal(s.T(), tc.expectedSavedCreditCard.PinCode, savedData[0].PinCode)
			assert.Equal(s.T(), tc.expectedSavedCreditCard.Metadata, savedData[0].Metadata)
		})
	}
}

func (s *CreditCardRepositoryTestSuite) insertTestUser(ownerID, login, password string) {
	query := "insert into infokeeper.user (id, login, password) values ($1, $2, $3)"
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(s.T(), err)

	_, err = s.pool.DB.Exec(context.Background(), query, ownerID, login, pass)
	assert.NoError(s.T(), err)
}
