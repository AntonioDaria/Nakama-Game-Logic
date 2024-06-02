package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"game_logic/src/mocks"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/heroiclabs/nakama-common/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRpcFileHandler(t *testing.T) {
	logger := &mocks.MockLogger{}

	mockNakama := new(mocks.CustomMockNakamaModule)
	db := &sql.DB{}

	// Set up file content and hash
	content := `{
    "example_data": "My data",
    "number": 1234567890
}`

	// Compact the JSON content
	var compactedContent bytes.Buffer
	if err := json.Compact(&compactedContent, []byte(content)); err != nil {
		log.Fatalf("Failed to compact JSON content: %v", err)
	}
	compactContent := compactedContent.Bytes()

	hash := md5.Sum(compactContent)
	calculatedHash := hex.EncodeToString(hash[:])

	// Create test file in a temporary directory
	tempDir := os.TempDir()
	testDir := filepath.Join(tempDir, "nakama/data/json_files")
	baseDir = testDir // Set the baseDir global variable for the tests
	err := os.MkdirAll(testDir, os.ModePerm)
	assert.NoError(t, err)

	testFilePath := filepath.Join(testDir, "core", "1.0.0.json")
	err = os.MkdirAll(filepath.Dir(testFilePath), os.ModePerm)
	assert.NoError(t, err)

	err = os.WriteFile(testFilePath, compactContent, os.ModePerm)
	assert.NoError(t, err)

	// Verify file existence
	if _, err := os.Stat(testFilePath); os.IsNotExist(err) {
		log.Fatalf("Test file does not exist: %s", testFilePath)
	}
	tests := []struct {
		name            string
		payload         string
		expectedType    string
		expectedVersion string
		expectedHash    string
		expectedContent json.RawMessage
		expectedError   error
	}{
		{
			name: "Valid payload with correct hash",
			payload: fmt.Sprintf(`{
				"type": "core",
				"version": "1.0.0",
				"hash": "%s"
			}`, calculatedHash),
			expectedType:    "core",
			expectedVersion: "1.0.0",
			expectedHash:    calculatedHash,
			expectedContent: compactContent,
			expectedError:   nil,
		},
		{
			name: "File doesn't exist, return error",
			payload: `{
				"type": "core",
				"version": "2.0.0",
				"hash": "wronghash"
			}`,
			expectedError: errFileNotFound,
		},
		{
			name: "Defaults are used if not present in payload",
			payload: fmt.Sprintf(`{
				"hash": "%s"
			}`, calculatedHash),
			expectedType:    "core",
			expectedVersion: "1.0.0",
			expectedHash:    calculatedHash,
			expectedContent: compactContent,
			expectedError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockNakama.On("StorageWrite", mock.Anything, mock.Anything).Return([]*api.StorageObjectAck{}, nil)

			// Call RpcFileHandler
			responseJson, err := RpcFileHandler(context.Background(), logger, db, mockNakama, tt.payload)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)

				// Parse response
				var response Response
				err = json.Unmarshal([]byte(responseJson), &response)
				assert.NoError(t, err)

				// Validate response fields
				assert.Equal(t, tt.expectedType, response.Type)
				assert.Equal(t, tt.expectedVersion, response.Version)
				assert.Equal(t, tt.expectedHash, response.Hash)
				assert.Equal(t, tt.expectedContent, response.Content)
			}
		})
	}
}

func TestRpcFileHandler_StatusCodes(t *testing.T) {
	logger := &mocks.MockLogger{}
	mockNakama := new(mocks.CustomMockNakamaModule)
	db := &sql.DB{}

	tests := []struct {
		name          string
		payload       string
		expectedError error
	}{
		{
			name:          "Invalid payload",
			payload:       "",
			expectedError: errBadInput,
		},
		{
			name: "File doesn't exist, return error",
			payload: `{
				"type": "core",
				"version": "2.0.0",
				"hash": "wronghash"
			}`,
			expectedError: errFileNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockNakama.On("StorageWrite", mock.Anything, mock.Anything).Return([]*api.StorageObjectAck{}, nil)

			// Call RpcFileHandler
			_, err := RpcFileHandler(context.Background(), logger, db, mockNakama, tt.payload)

			// Check errors
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
