package handler_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/carry"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/carry"
)

// SimpleTestLogger for testing - doesn't require mock expectations
type SimpleTestLogger struct{}

func (l *SimpleTestLogger) Debug(ctx context.Context, service, message string, metadata ...map[string]interface{}) {
}
func (l *SimpleTestLogger) Info(ctx context.Context, service, message string, metadata ...map[string]interface{}) {
}
func (l *SimpleTestLogger) Warning(ctx context.Context, service, message string, metadata ...map[string]interface{}) {
}
func (l *SimpleTestLogger) Error(ctx context.Context, service, message string, err error, metadata ...map[string]interface{}) {
}
func (l *SimpleTestLogger) Fatal(ctx context.Context, service, message string, err error, metadata ...map[string]interface{}) {
}
func (l *SimpleTestLogger) LogRequest(ctx context.Context, service, method, endpoint, userAgent, ipAddress string, userID *int, requestID string) {
}
func (l *SimpleTestLogger) LogResponse(ctx context.Context, service, requestID string, statusCode, executionTimeMs int, metadata ...map[string]interface{}) {
}

func TestCarryHandler_SetLogger(t *testing.T) {
	// arrange
	mockService := &mocks.CarryServiceMock{}
	testLogger := &SimpleTestLogger{}

	carryHandler := handler.NewCarryHandler(mockService)

	// act
	carryHandler.SetLogger(testLogger)

	// assert - no error should occur
	require.NotNil(t, carryHandler)
}

func TestCarryHandler_SetLogger_NilLogger(t *testing.T) {
	// arrange
	mockService := &mocks.CarryServiceMock{}

	carryHandler := handler.NewCarryHandler(mockService)

	// act
	carryHandler.SetLogger(nil)

	// assert - no error should occur with nil logger
	require.NotNil(t, carryHandler)
}
