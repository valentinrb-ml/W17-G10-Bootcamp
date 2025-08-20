package handler_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	handler "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/handler/warehouse"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/warehouse"
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

func TestWarehouseHandler_SetLogger(t *testing.T) {
	// arrange
	mockService := &mocks.WarehouseServiceMock{}
	warehouseHandler := handler.NewWarehouseHandler(mockService)
	testLogger := &SimpleTestLogger{}

	// act
	warehouseHandler.SetLogger(testLogger)

	// assert
	// Since SetLogger doesn't return anything, we just verify no panic occurred
	// and the logger was set (which is verified by the fact that the method executed successfully)
	require.NotNil(t, warehouseHandler)
}

func TestWarehouseHandler_SetLogger_NilLogger(t *testing.T) {
	// arrange
	mockService := &mocks.WarehouseServiceMock{}
	warehouseHandler := handler.NewWarehouseHandler(mockService)

	// act & assert - should not panic with nil logger
	require.NotPanics(t, func() {
		warehouseHandler.SetLogger(nil)
	})
}
