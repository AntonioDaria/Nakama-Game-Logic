package mocks

import (
	"context"
	"log"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/stretchr/testify/mock"
)

// Mock Nakama module
type CustomMockNakamaModule struct {
	runtime.NakamaModule
	mock.Mock
}

func (m *CustomMockNakamaModule) StorageWrite(ctx context.Context, writes []*runtime.StorageWrite) ([]*api.StorageObjectAck, error) {
	args := m.Called(ctx, writes)
	return args.Get(0).([]*api.StorageObjectAck), args.Error(1)
}

// Mock Logger
type MockLogger struct{}

func (l *MockLogger) Debug(format string, v ...interface{}) {
	log.Printf(format, v...)
}
func (l *MockLogger) Info(format string, v ...interface{})  {}
func (l *MockLogger) Warn(format string, v ...interface{})  {}
func (l *MockLogger) Error(format string, v ...interface{}) {}
func (l *MockLogger) Fatal(format string, v ...interface{}) {}

func (l *MockLogger) Fields() map[string]interface{} {
	return nil
}
func (l *MockLogger) WithField(key string, value interface{}) runtime.Logger {
	return nil
}
func (l *MockLogger) WithFields(map[string]interface{}) runtime.Logger {
	return nil
}
