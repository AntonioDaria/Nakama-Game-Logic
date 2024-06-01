package main

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

// HealthCheckResponse is the response object for the healthcheck RPC function.
type HealthCheckResponse struct {
	Health string `json:"health"`
}

// HealthCheck is an RPC function that returns a healthcheck response.
func RpcHealthCheck(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	logger.Debug("HealthCheck function was invoked")
	response := HealthCheckResponse{
		Health: "OK",
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		logger.Error("Error encoding response: %v", err)
		return "", err
	}
	return string(responseJSON), nil
}
