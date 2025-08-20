package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/warehouse"
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

func TestWarehouseService_SetLogger(t *testing.T) {
	// arrange
	mockRepo := &mocks.WarehouseRepositoryMock{}
	warehouseService := service.NewWarehouseService(mockRepo)
	testLogger := &SimpleTestLogger{}

	// act
	warehouseService.SetLogger(testLogger)

	// assert
	// Since SetLogger doesn't return anything, we just verify no panic occurred
	require.NotNil(t, warehouseService)
}

func TestWarehouseService_SetLogger_NilLogger(t *testing.T) {
	// arrange
	mockRepo := &mocks.WarehouseRepositoryMock{}
	warehouseService := service.NewWarehouseService(mockRepo)

	// act & assert - should not panic with nil logger
	require.NotPanics(t, func() {
		warehouseService.SetLogger(nil)
	})
}
