package logger

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
)

// LogLevel represents the different logging levels
type LogLevel string

const (
	DEBUG   LogLevel = "DEBUG"
	INFO    LogLevel = "INFO"
	WARNING LogLevel = "WARNING"
	ERROR   LogLevel = "ERROR"
	FATAL   LogLevel = "FATAL"
)

// Logger interface defines the basic methods for logging
type Logger interface {
	Debug(ctx context.Context, service, message string, metadata ...map[string]interface{})
	Info(ctx context.Context, service, message string, metadata ...map[string]interface{})
	Warning(ctx context.Context, service, message string, metadata ...map[string]interface{})
	Error(ctx context.Context, service, message string, err error, metadata ...map[string]interface{})
	Fatal(ctx context.Context, service, message string, err error, metadata ...map[string]interface{})
	LogRequest(ctx context.Context, service, method, endpoint, userAgent, ipAddress string, userID *int, requestID string)
	LogResponse(ctx context.Context, service, requestID string, statusCode, executionTimeMs int, metadata ...map[string]interface{})
}

// DatabaseLogger implements Logger using database
type DatabaseLogger struct {
	db *sql.DB
}

// NewDatabaseLogger creates a new instance of the database logger
func NewDatabaseLogger(db *sql.DB) Logger {
	return &DatabaseLogger{db: db}
}

// logEntry saves a log entry to the database
func (l *DatabaseLogger) logEntry(ctx context.Context, level LogLevel, service, message string, errorDetails string, metadata map[string]interface{}) {
	// Get information from context
	requestID := l.getRequestID(ctx)
	endpoint := l.getEndpoint(ctx)
	method := l.getMethod(ctx)
	userID := l.getUserID(ctx)
	ipAddress := l.getIPAddress(ctx)
	userAgent := l.getUserAgent(ctx)
	statusCode := l.getStatusCode(ctx)
	executionTime := l.getExecutionTime(ctx)

	// Add caller information to metadata
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	// Get caller information
	if pc, file, line, ok := runtime.Caller(3); ok {
		parts := strings.Split(file, "/")
		filename := parts[len(parts)-1]

		// Add multi-level context for better clarity
		var layerContext string
		if len(parts) >= 3 {
			// Take the last 3 levels: repository/warehouse/warehouse.go
			parentDir := parts[len(parts)-3] // repository, service, handler
			moduleDir := parts[len(parts)-2] // warehouse
			layerContext = parentDir + "/" + moduleDir + "/" + filename
		} else if len(parts) >= 2 {
			// Fallback to 2 levels if not enough parts
			parentDir := parts[len(parts)-2]
			layerContext = parentDir + "/" + filename
		} else {
			// Fallback to filename only
			layerContext = filename
		}

		funcName := runtime.FuncForPC(pc).Name()
		if idx := strings.LastIndex(funcName, "."); idx != -1 {
			funcName = funcName[idx+1:]
		}

		metadata["caller_file"] = layerContext // Now: repository/warehouse/warehouse.go
		metadata["caller_line"] = line
		metadata["caller_function"] = funcName

		// Also add filename only for backwards compatibility
		metadata["filename"] = filename

		// Add specific layer for easy filtering
		if len(parts) >= 3 {
			metadata["layer"] = parts[len(parts)-3] // repository, service, handler
		}
	}

	// Serialize metadata to JSON
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		metadataJSON = []byte("{}")
	}

	query := `
		INSERT INTO logs (
			level, service, endpoint, method, user_id, request_id, 
			message, metadata, execution_time_ms, status_code, 
			error_details, ip_address, user_agent
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = l.db.ExecContext(ctx, query,
		string(level), service, endpoint, method, userID, requestID,
		message, string(metadataJSON), executionTime, statusCode,
		errorDetails, ipAddress, userAgent,
	)
	if err != nil {
		// In case of error saving to DB, print to console as fallback
		fmt.Printf("ERROR: Failed to save log to database: %v. Original log: [%s] %s - %s\n", err, level, service, message)
	}
}

// Helper methods to extract information from context
func (l *DatabaseLogger) getRequestID(ctx context.Context) string {
	if reqID := ctx.Value("request_id"); reqID != nil {
		if id, ok := reqID.(string); ok {
			return id
		}
	}
	return ""
}

func (l *DatabaseLogger) getEndpoint(ctx context.Context) *string {
	if endpoint := ctx.Value("endpoint"); endpoint != nil {
		if ep, ok := endpoint.(string); ok {
			return &ep
		}
	}
	return nil
}

func (l *DatabaseLogger) getMethod(ctx context.Context) *string {
	if method := ctx.Value("method"); method != nil {
		if m, ok := method.(string); ok {
			return &m
		}
	}
	return nil
}

func (l *DatabaseLogger) getUserID(ctx context.Context) *int {
	if userID := ctx.Value("user_id"); userID != nil {
		if id, ok := userID.(int); ok {
			return &id
		}
	}
	return nil
}

func (l *DatabaseLogger) getIPAddress(ctx context.Context) *string {
	if ip := ctx.Value("ip_address"); ip != nil {
		if ipAddr, ok := ip.(string); ok {
			return &ipAddr
		}
	}
	return nil
}

func (l *DatabaseLogger) getUserAgent(ctx context.Context) *string {
	if ua := ctx.Value("user_agent"); ua != nil {
		if userAgent, ok := ua.(string); ok {
			return &userAgent
		}
	}
	return nil
}

func (l *DatabaseLogger) getStatusCode(ctx context.Context) *int {
	if sc := ctx.Value("status_code"); sc != nil {
		if code, ok := sc.(int); ok {
			return &code
		}
	}
	return nil
}

func (l *DatabaseLogger) getExecutionTime(ctx context.Context) *int {
	if et := ctx.Value("execution_time_ms"); et != nil {
		if time, ok := et.(int); ok {
			return &time
		}
	}
	return nil
}

// Debug logs a debug message
func (l *DatabaseLogger) Debug(ctx context.Context, service, message string, metadata ...map[string]interface{}) {
	var meta map[string]interface{}
	if len(metadata) > 0 {
		meta = metadata[0]
	}
	l.logEntry(ctx, DEBUG, service, message, "", meta)
}

// Info logs an informational message
func (l *DatabaseLogger) Info(ctx context.Context, service, message string, metadata ...map[string]interface{}) {
	var meta map[string]interface{}
	if len(metadata) > 0 {
		meta = metadata[0]
	}
	l.logEntry(ctx, INFO, service, message, "", meta)
}

// Warning logs a warning message
func (l *DatabaseLogger) Warning(ctx context.Context, service, message string, metadata ...map[string]interface{}) {
	var meta map[string]interface{}
	if len(metadata) > 0 {
		meta = metadata[0]
	}
	l.logEntry(ctx, WARNING, service, message, "", meta)
}

// Error logs an error message
func (l *DatabaseLogger) Error(ctx context.Context, service, message string, err error, metadata ...map[string]interface{}) {
	var meta map[string]interface{}
	if len(metadata) > 0 {
		meta = metadata[0]
	}

	errorDetails := ""
	if err != nil {
		errorDetails = err.Error()
	}

	l.logEntry(ctx, ERROR, service, message, errorDetails, meta)
}

// Fatal logs a fatal error message
func (l *DatabaseLogger) Fatal(ctx context.Context, service, message string, err error, metadata ...map[string]interface{}) {
	var meta map[string]interface{}
	if len(metadata) > 0 {
		meta = metadata[0]
	}

	errorDetails := ""
	if err != nil {
		errorDetails = err.Error()
	}

	l.logEntry(ctx, FATAL, service, message, errorDetails, meta)
}

// LogRequest logs the start of an HTTP request
func (l *DatabaseLogger) LogRequest(ctx context.Context, service, method, endpoint, userAgent, ipAddress string, userID *int, requestID string) {
	metadata := map[string]interface{}{
		"method":     method,
		"endpoint":   endpoint,
		"user_agent": userAgent,
		"ip_address": ipAddress,
		"request_id": requestID,
	}

	if userID != nil {
		metadata["user_id"] = *userID
	}

	query := `
		INSERT INTO logs (
			level, service, endpoint, method, user_id, request_id, 
			message, metadata, ip_address, user_agent
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		metadataJSON = []byte("{}")
	}

	message := fmt.Sprintf("Request started: %s %s", method, endpoint)

	_, err = l.db.ExecContext(ctx, query,
		string(INFO), service, endpoint, method, userID, requestID,
		message, string(metadataJSON), ipAddress, userAgent,
	)
	if err != nil {
		fmt.Printf("ERROR: Failed to save request log to database: %v\n", err)
	}
}

// LogResponse logs the end of an HTTP request with its response
func (l *DatabaseLogger) LogResponse(ctx context.Context, service, requestID string, statusCode, executionTimeMs int, metadata ...map[string]interface{}) {
	var meta map[string]interface{}
	if len(metadata) > 0 {
		meta = metadata[0]
	} else {
		meta = make(map[string]interface{})
	}

	meta["status_code"] = statusCode
	meta["execution_time_ms"] = executionTimeMs
	meta["request_id"] = requestID

	// Get method and endpoint from context if available
	endpoint := l.getEndpoint(ctx)
	method := l.getMethod(ctx)

	var endpointStr, methodStr string
	if endpoint != nil {
		endpointStr = *endpoint
	}
	if method != nil {
		methodStr = *method
	}

	message := fmt.Sprintf("Request completed: %s %s (%d, %dms)", methodStr, endpointStr, statusCode, executionTimeMs)

	metadataJSON, err := json.Marshal(meta)
	if err != nil {
		metadataJSON = []byte("{}")
	}

	query := `
		INSERT INTO logs (
			level, service, request_id, message, metadata, 
			execution_time_ms, status_code, endpoint, method
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = l.db.ExecContext(ctx, query,
		string(INFO), service, requestID, message, string(metadataJSON),
		executionTimeMs, statusCode, endpointStr, methodStr,
	)
	if err != nil {
		fmt.Printf("ERROR: Failed to save response log to database: %v\n", err)
	}
}
