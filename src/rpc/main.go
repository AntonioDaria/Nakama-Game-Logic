package main

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	logger.Info("Module loaded")
	// Register the RPC function for health checks
	if err := initializer.RegisterRpc("healthCheck", RpcHealthCheck); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	// Set the base directory for JSON files
	baseDir = "/nakama/data/json_files"

	// Register the RPC function for file handling
	if err := initializer.RegisterRpc("fileHandler", RpcFileHandler); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}
	return nil
}
