package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	repository "github.com/varobledo_meli/W17-G10-Bootcamp.git/internal/repository/carry"
	"github.com/varobledo_meli/W17-G10-Bootcamp.git/testhelpers"
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

func TestCarryRepository_SetLogger(t *testing.T) {
	// arrange
	mock, db := testhelpers.CreateMockDB()
	defer db.Close()

	testLogger := &SimpleTestLogger{}

	carryRepo := repository.NewCarryRepository(db)

	// act
	carryRepo.SetLogger(testLogger)

	// assert - no error should occur
	require.NotNil(t, carryRepo)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCarryRepository_SetLogger_NilLogger(t *testing.T) {
	// arrange
	mock, db := testhelpers.CreateMockDB()
	defer db.Close()

	carryRepo := repository.NewCarryRepository(db)

	// act
	carryRepo.SetLogger(nil)

	// assert - no error should occur with nil logger
	require.NotNil(t, carryRepo)
	require.NoError(t, mock.ExpectationsWereMet())
}
