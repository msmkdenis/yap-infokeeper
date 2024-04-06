package repository

import (
	"context"
	_ "embed"
	"encoding/json"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/msmkdenis/yap-infokeeper/internal/model"
	"github.com/msmkdenis/yap-infokeeper/internal/text_data/specification"
)

//go:embed fixtures/text_data.json
var fixtureTextData string

func (s *TextDataRepositoryTestSuite) Test_SelectAll() {
	textData := s.prepareTextData()
	s.insertTestUser(textData[0].OwnerID, "login 1", "password 1")
	s.insertTestUser(textData[1].OwnerID, "login 2", "password 2")
	s.insertTestData(textData)

	testCases := []struct {
		name          string
		spec          specification.TextDataSpecification
		expectedError error
		filterFunc    func(data model.TextData) bool
	}{
		{
			name: "Success - only OwnerID",
			spec: specification.TextDataSpecification{
				OwnerID: textData[0].OwnerID,
			},
			expectedError: nil,
			filterFunc: func(data model.TextData) bool {
				return data.OwnerID == textData[0].OwnerID
			},
		},
		{
			name: "Success - ownerID, data",
			spec: specification.TextDataSpecification{
				OwnerID: textData[1].OwnerID,
				Data:    "is",
			},
			expectedError: nil,
			filterFunc: func(data model.TextData) bool {
				return data.OwnerID == textData[1].OwnerID && strings.Contains(data.Data, "is")
			},
		},
		{
			name: "Success - ownerID, data, metadata",
			spec: specification.TextDataSpecification{
				OwnerID:  textData[1].OwnerID,
				Data:     "is",
				Metadata: "relationships",
			},
			expectedError: nil,
			filterFunc: func(data model.TextData) bool {
				return data.OwnerID == textData[1].OwnerID &&
					strings.Contains(data.Data, "is") &&
					strings.Contains(data.Metadata, "relationships")
			},
		},
		{
			name: "Success - ownerID, createdAfter, createdBefore",
			spec: specification.TextDataSpecification{
				OwnerID:       textData[1].OwnerID,
				CreatedAfter:  time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				CreatedBefore: time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedError: nil,
			filterFunc: func(data model.TextData) bool {
				if data.OwnerID == textData[1].OwnerID &&
					data.CreatedAt.After(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)) &&
					data.CreatedAt.Before(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)) {
					return true
				}
				return false
			},
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			selectedData, err := s.textDataRepository.SelectAll(context.Background(), &tc.spec)
			assert.Equal(s.T(), tc.expectedError, err)

			expectedData := s.filterTextData(textData, tc.filterFunc)
			assert.Equal(s.T(), len(expectedData), len(selectedData))
			assert.Equal(s.T(), expectedData, selectedData)
		})
	}
}

func (s *TextDataRepositoryTestSuite) insertTestData(textData []model.TextData) {
	rows := make([][]interface{}, len(textData))
	for i, text := range textData {
		row := []interface{}{text.ID, text.OwnerID, text.Data, text.CreatedAt, text.Metadata}
		rows[i] = row
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := s.pool.DB.CopyFrom(
		ctx,
		pgx.Identifier{"infokeeper", "text_data"},
		[]string{"id", "owner_id", "data", "created_at", "metadata"},
		pgx.CopyFromRows(rows),
	)
	assert.Equal(s.T(), int64(len(textData)), count)
	assert.NoError(s.T(), err)
}

func (s *TextDataRepositoryTestSuite) prepareTextData() []model.TextData {
	var textData []model.TextData
	err := json.NewDecoder(strings.NewReader(fixtureTextData)).Decode(&textData)
	require.NoError(s.T(), err)
	return textData
}

func (s *TextDataRepositoryTestSuite) filterTextData(textData []model.TextData, fn func(txt model.TextData) bool) []model.TextData {
	filteredTextData := make([]model.TextData, 0)
	for _, text := range textData {
		if fn(text) {
			filteredTextData = append(filteredTextData, text)
		}
	}
	return filteredTextData
}
