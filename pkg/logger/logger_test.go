package logger

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestDatabaseLogger_InfoWithMetadata(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := NewDatabaseLogger(db)
	ctx := context.Background()

	metadata := map[string]interface{}{
		"user_id": 123,
		"action":  "create",
	}

	mock.ExpectExec("INSERT INTO logs").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	logger.Info(ctx, "test-service", "Test message", metadata)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDatabaseLogger_ErrorWithMetadata(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := NewDatabaseLogger(db)
	ctx := context.Background()

	testErr := assert.AnError
	metadata := map[string]interface{}{
		"request_id": "123",
	}

	mock.ExpectExec("INSERT INTO logs").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	logger.Error(ctx, "test-service", "Error message", testErr, metadata)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDatabaseLogger_LogRequestWithNilUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := NewDatabaseLogger(db)
	ctx := context.Background()

	mock.ExpectExec("INSERT INTO logs").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	logger.LogRequest(ctx, "test-service", "POST", "/test", "test-agent", "192.168.1.1", nil, "req-456")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDatabaseLogger_LogResponseWithMetadata(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := NewDatabaseLogger(db)
	ctx := context.Background()

	metadata := map[string]interface{}{
		"response_size": 1024,
	}

	mock.ExpectExec("INSERT INTO logs").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	logger.LogResponse(ctx, "test-service", "req-789", 201, 250, metadata)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDatabaseLogger_WithContextValues(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := NewDatabaseLogger(db)

	// Test with various context values
	ctx := context.WithValue(context.Background(), "endpoint", "/api/test")
	ctx = context.WithValue(ctx, "method", "GET")
	ctx = context.WithValue(ctx, "user_id", 999)
	ctx = context.WithValue(ctx, "request_id", "req-ctx-123")

	mock.ExpectExec("INSERT INTO logs").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	logger.Info(ctx, "test-service", "Context test message")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDatabaseLogger_DatabaseErrorHandling(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := NewDatabaseLogger(db)
	ctx := context.Background()

	// Test database error handling
	mock.ExpectExec("INSERT INTO logs").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(assert.AnError)

	// Should not panic even with database error
	require.NotPanics(t, func() {
		logger.Info(ctx, "test-service", "Error test message")
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDatabaseLogger_WarningWithMetadata(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := NewDatabaseLogger(db)
	ctx := context.Background()

	metadata := map[string]interface{}{
		"warning_code": "W001",
	}

	mock.ExpectExec("INSERT INTO logs").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	logger.Warning(ctx, "test-service", "Warning message", metadata)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDatabaseLogger_DebugWithMetadata(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := NewDatabaseLogger(db)
	ctx := context.Background()

	metadata := map[string]interface{}{
		"debug_info": "detailed debug data",
	}

	mock.ExpectExec("INSERT INTO logs").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	logger.Debug(ctx, "test-service", "Debug message", metadata)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDatabaseLogger_FatalWithMetadata(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := NewDatabaseLogger(db)
	ctx := context.Background()

	testErr := assert.AnError
	metadata := map[string]interface{}{
		"crash_dump": "system_state",
	}

	mock.ExpectExec("INSERT INTO logs").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	logger.Fatal(ctx, "test-service", "Fatal message", testErr, metadata)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDatabaseLogger_AllContextValues(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := NewDatabaseLogger(db)

	// Test with all possible context values
	ctx := context.WithValue(context.Background(), "endpoint", "/comprehensive/test")
	ctx = context.WithValue(ctx, "method", "PUT")
	ctx = context.WithValue(ctx, "user_id", 555)
	ctx = context.WithValue(ctx, "request_id", "req-comprehensive-789")
	ctx = context.WithValue(ctx, "ip_address", "10.0.0.1")
	ctx = context.WithValue(ctx, "user_agent", "comprehensive-test-agent")
	ctx = context.WithValue(ctx, "status_code", 201)
	ctx = context.WithValue(ctx, "execution_time", 300)

	mock.ExpectExec("INSERT INTO logs").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	logger.Info(ctx, "test-service", "Comprehensive context test")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDatabaseLogger_MultipleMetadataArgs(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := NewDatabaseLogger(db)
	ctx := context.Background()

	metadata1 := map[string]interface{}{
		"first": "value1",
	}
	metadata2 := map[string]interface{}{
		"second": "value2",
	}

	mock.ExpectExec("INSERT INTO logs").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Test multiple metadata arguments
	logger.Info(ctx, "test-service", "Multiple metadata test", metadata1, metadata2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDatabaseLogger_LogRequestDatabase_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := NewDatabaseLogger(db)
	ctx := context.Background()

	userID := 123

	mock.ExpectExec("INSERT INTO logs").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(assert.AnError)

	// Should not panic even with database error
	require.NotPanics(t, func() {
		logger.LogRequest(ctx, "test-service", "GET", "/test", "test-agent", "127.0.0.1", &userID, "req-123")
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDatabaseLogger_LogResponseDatabase_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := NewDatabaseLogger(db)
	ctx := context.Background()

	mock.ExpectExec("INSERT INTO logs").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(assert.AnError)

	// Should not panic even with database error
	require.NotPanics(t, func() {
		logger.LogResponse(ctx, "test-service", "req-123", 200, 150)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}
