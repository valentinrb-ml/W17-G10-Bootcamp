package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	service "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/service/carry"
	mocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/carry"
	geographyMocks "github.com/varobledo_meli/W17-G10-Bootcamp.git/mocks/geography"
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

func TestCarryService_SetLogger(t *testing.T) {
	// arrange
	mockRepo := &mocks.CarryRepositoryMock{}
	mockGeography := &geographyMocks.GeographyRepositoryMock{}
	testLogger := &SimpleTestLogger{}

	carryService := service.NewCarryService(mockRepo, mockGeography)

	// act
	carryService.SetLogger(testLogger)

	// assert - no error should occur
	require.NotNil(t, carryService)
}

func TestCarryService_SetLogger_NilLogger(t *testing.T) {
	// arrange
	mockRepo := &mocks.CarryRepositoryMock{}
	mockGeography := &geographyMocks.GeographyRepositoryMock{}

	carryService := service.NewCarryService(mockRepo, mockGeography)

	// act
	carryService.SetLogger(nil)

	// assert - no error should occur with nil logger
	require.NotNil(t, carryService)
}
