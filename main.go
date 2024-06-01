package main

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	logger.Info("Module loaded")
	// Register the RPC function
	if err := initializer.RegisterRpc("healthCheck", RpcHealthCheck); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}
	return nil
}
