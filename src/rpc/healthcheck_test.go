package main

import (
	"context"
	"database/sql"
	"testing"

	"game_logic/src/mocks"

	"github.com/heroiclabs/nakama-common/runtime"
)

func TestRpcHealthCheck(t *testing.T) {
	type args struct {
		ctx     context.Context
		logger  runtime.Logger
		db      *sql.DB
		nk      runtime.NakamaModule
		payload string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test health check",
			args: args{
				ctx:     context.Background(),
				logger:  &mocks.MockLogger{},
				db:      &sql.DB{},
				nk:      &mocks.CustomMockNakamaModule{},
				payload: "",
			},
			want:    `{"health":"OK"}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RpcHealthCheck(tt.args.ctx, tt.args.logger, tt.args.db, tt.args.nk, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("RpcHealthCheck() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RpcHealthCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
