package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/heroiclabs/nakama-common/runtime"
)

type Payload struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Hash    string `json:"hash"`
}

type Response struct {
	Type    string          `json:"type"`
	Version string          `json:"version"`
	Hash    string          `json:"hash"`
	Content json.RawMessage `json:"content"`
}

var baseDir string

const (
	INVALID_REQUEST = 3
	NOT_FOUND       = 5
	INTERNAL        = 13
)

var (
	errBadInput      = runtime.NewError("input contained invalid data", INVALID_REQUEST)
	errFileNotFound  = runtime.NewError("file not found", NOT_FOUND)
	errInternalError = runtime.NewError("internal server error", INTERNAL)
)

func RpcFileHandler(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("RpcFileHandler function was invoked with payload: %s", payload)

	// Set defaults
	req := Payload{
		Type:    "core",
		Version: "1.0.0",
		Hash:    "",
	}

	if payload == "" {
		logger.Error("Payload is empty")
		return "", errBadInput

	}

	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		logger.Error("Failed to unmarshal payload: %v", err)
		return "", errBadInput
	}

	// Use defaults if values are not provided in the payload
	if req.Type == "" {
		req.Type = "core"
	}
	if req.Version == "" {
		req.Version = "1.0.0"
	}
	if req.Hash == "" {
		req.Hash = ""
	}

	// Read file
	filePath := fmt.Sprintf("%s/%s/%s.json", baseDir, req.Type, req.Version)
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		logger.Error("Failed to read file at %s: %v", filePath, err)
		return "", errFileNotFound
	}
	logger.Debug("Read file content: %s", string(fileContent))

	// Compact the JSON content
	var compactedContent bytes.Buffer
	if err := json.Compact(&compactedContent, fileContent); err != nil {
		logger.Error("Failed to compact JSON content: %v", err)
		return "", err
	}
	compactFileContent := compactedContent.Bytes()

	// Calculate hash
	hash := md5.Sum(compactFileContent)
	calculatedHash := hex.EncodeToString(hash[:])
	logger.Debug("Calculated hash: %s", calculatedHash)

	// // Prepare response
	response := Response{
		Type:    req.Type,
		Version: req.Version,
		Hash:    calculatedHash,
		Content: nil,
	}

	// Check hash and save content
	if req.Hash != "" && req.Hash == calculatedHash {

		response.Content = json.RawMessage(compactFileContent)

		// Save to Nakama's storage engine
		storageWrite := []*runtime.StorageWrite{
			{
				Collection:      req.Type,
				Key:             req.Version,
				Value:           string(compactFileContent),
				PermissionRead:  2, // Public read
				PermissionWrite: 0, // Owner write
			},
		}
		if _, err := nk.StorageWrite(ctx, storageWrite); err != nil {
			logger.Error("Failed to write to storage: %v", err)
			return "", errInternalError
		}
		logger.Debug("Content saved to storage")
	} else {
		response.Content = nil
		logger.Debug("Hashes do not match or no hash provided")
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		logger.Error("Failed to marshal response: %v", err)
		return "", errInternalError
	}

	return string(responseJson), nil
}
