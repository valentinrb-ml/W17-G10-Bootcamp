package logger

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLogger for testing
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(ctx context.Context, service, message string, metadata ...map[string]interface{}) {
	args := []interface{}{ctx, service, message}
	for _, meta := range metadata {
		args = append(args, meta)
	}
	m.Called(args...)
}

func (m *MockLogger) Info(ctx context.Context, service, message string, metadata ...map[string]interface{}) {
	args := []interface{}{ctx, service, message}
	for _, meta := range metadata {
		args = append(args, meta)
	}
	m.Called(args...)
}

func (m *MockLogger) Warning(ctx context.Context, service, message string, metadata ...map[string]interface{}) {
	args := []interface{}{ctx, service, message}
	for _, meta := range metadata {
		args = append(args, meta)
	}
	m.Called(args...)
}

func (m *MockLogger) Error(ctx context.Context, service, message string, err error, metadata ...map[string]interface{}) {
	args := []interface{}{ctx, service, message, err}
	for _, meta := range metadata {
		args = append(args, meta)
	}
	m.Called(args...)
}

func (m *MockLogger) Fatal(ctx context.Context, service, message string, err error, metadata ...map[string]interface{}) {
	args := []interface{}{ctx, service, message, err}
	for _, meta := range metadata {
		args = append(args, meta)
	}
	m.Called(args...)
}

func (m *MockLogger) LogRequest(ctx context.Context, service, method, endpoint, userAgent, ipAddress string, userID *int, requestID string) {
	m.Called(ctx, service, method, endpoint, userAgent, ipAddress, userID, requestID)
}

func (m *MockLogger) LogResponse(ctx context.Context, service, requestID string, statusCode, executionTimeMs int, metadata ...map[string]interface{}) {
	args := []interface{}{ctx, service, requestID, statusCode, executionTimeMs}
	for _, meta := range metadata {
		args = append(args, meta)
	}
	m.Called(args...)
}

// Test handler that extracts context values for verification
func testHandler(t *testing.T, expectedValues map[string]interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Verify context values
		for key, expectedValue := range expectedValues {
			actualValue := r.Context().Value(key)
			assert.Equal(t, expectedValue, actualValue, "Context value for key %s should match", key)
		}

		// Verify Request ID is a valid UUID
		requestID := r.Context().Value("request_id")
		if requestID != nil {
			_, err := uuid.Parse(requestID.(string))
			assert.NoError(t, err, "Request ID should be a valid UUID")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}

func TestLoggingMiddleware_RequestIDGeneration(t *testing.T) {
	mockLogger := new(MockLogger)
	middleware := LoggingMiddleware(mockLogger)

	// Mock expectations
	mockLogger.On("LogRequest", mock.Anything, mock.AnythingOfType("string"), "GET", "/test", mock.AnythingOfType("string"), mock.AnythingOfType("string"), (*int)(nil), mock.AnythingOfType("string")).Return()
	mockLogger.On("LogResponse", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), 200, mock.AnythingOfType("int"), mock.Anything).Return()

	handler := middleware(testHandler(t, map[string]interface{}{}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Verify Request ID header is set
	requestID := w.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID, "X-Request-ID header should be set")

	// Verify it's a valid UUID
	_, err := uuid.Parse(requestID)
	assert.NoError(t, err, "X-Request-ID should be a valid UUID")

	mockLogger.AssertExpectations(t)
}

func TestLoggingMiddleware_RequestIDUniqueness(t *testing.T) {
	mockLogger := new(MockLogger)
	middleware := LoggingMiddleware(mockLogger)

	// Allow multiple calls
	mockLogger.On("LogRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	mockLogger.On("LogResponse", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	requestIDs := make(map[string]bool)

	// Generate multiple request IDs
	for i := 0; i < 100; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		requestID := w.Header().Get("X-Request-ID")
		assert.NotEmpty(t, requestID)

		// Verify uniqueness
		assert.False(t, requestIDs[requestID], "Request ID should be unique: %s", requestID)
		requestIDs[requestID] = true
	}
}

func TestLoggingMiddleware_ContextPropagation(t *testing.T) {
	mockLogger := new(MockLogger)
	middleware := LoggingMiddleware(mockLogger)

	expectedValues := map[string]interface{}{
		"endpoint":   "/warehouses",
		"method":     "POST",
		"ip_address": "192.168.1.1",
	}

	mockLogger.On("LogRequest", mock.Anything, "warehouse-service", "POST", "/warehouses", mock.AnythingOfType("string"), "192.168.1.1", (*int)(nil), mock.AnythingOfType("string")).Return()
	mockLogger.On("LogResponse", mock.Anything, "warehouse-service", mock.AnythingOfType("string"), 200, mock.AnythingOfType("int"), mock.Anything).Return()

	handler := middleware(testHandler(t, expectedValues))

	req := httptest.NewRequest("POST", "/warehouses", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	mockLogger.AssertExpectations(t)
}

func TestLoggingMiddleware_ResponseHeadersSet(t *testing.T) {
	mockLogger := new(MockLogger)
	middleware := LoggingMiddleware(mockLogger)

	mockLogger.On("LogRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	mockLogger.On("LogResponse", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created"))
	}))

	req := httptest.NewRequest("POST", "/warehouses", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Verify X-Request-ID header is set
	requestID := w.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID)

	// Verify it's a valid UUID
	_, err := uuid.Parse(requestID)
	assert.NoError(t, err)

	// Verify response
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "Created", w.Body.String())
}

func TestLoggingMiddleware_ExecutionTimeMeasurement(t *testing.T) {
	mockLogger := new(MockLogger)
	middleware := LoggingMiddleware(mockLogger)

	// Handler that takes some time
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))

	mockLogger.On("LogRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	mockLogger.On("LogResponse", mock.Anything, mock.Anything, mock.Anything, 200, mock.MatchedBy(func(execTime int) bool {
		return execTime >= 10 // Should be at least 10ms
	}), mock.Anything).Return()

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	mockLogger.AssertExpectations(t)
}

func TestLoggingMiddleware_StatusCodeCapture(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
	}{
		{"OK", http.StatusOK},
		{"Created", http.StatusCreated},
		{"Bad Request", http.StatusBadRequest},
		{"Not Found", http.StatusNotFound},
		{"Internal Server Error", http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockLogger := new(MockLogger)
			middleware := LoggingMiddleware(mockLogger)

			handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
			}))

			mockLogger.On("LogRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
			mockLogger.On("LogResponse", mock.Anything, mock.Anything, mock.Anything, tc.statusCode, mock.AnythingOfType("int"), mock.Anything).Return()

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tc.statusCode, w.Code)
			mockLogger.AssertExpectations(t)
		})
	}
}

func TestGetClientIP(t *testing.T) {
	testCases := []struct {
		name       string
		remoteAddr string
		headers    map[string]string
		expectedIP string
	}{
		{
			name:       "X-Forwarded-For single IP",
			remoteAddr: "10.0.0.1:12345",
			headers:    map[string]string{"X-Forwarded-For": "192.168.1.1"},
			expectedIP: "192.168.1.1",
		},
		{
			name:       "X-Forwarded-For multiple IPs",
			remoteAddr: "10.0.0.1:12345",
			headers:    map[string]string{"X-Forwarded-For": "192.168.1.1, 10.0.0.1, 172.16.0.1"},
			expectedIP: "192.168.1.1",
		},
		{
			name:       "X-Real-IP",
			remoteAddr: "10.0.0.1:12345",
			headers:    map[string]string{"X-Real-IP": "192.168.1.100"},
			expectedIP: "192.168.1.100",
		},
		{
			name:       "X-Forwarded-For takes precedence over X-Real-IP",
			remoteAddr: "10.0.0.1:12345",
			headers: map[string]string{
				"X-Forwarded-For": "192.168.1.1",
				"X-Real-IP":       "192.168.1.100",
			},
			expectedIP: "192.168.1.1",
		},
		{
			name:       "RemoteAddr fallback with port",
			remoteAddr: "192.168.1.1:12345",
			headers:    map[string]string{},
			expectedIP: "192.168.1.1",
		},
		{
			name:       "RemoteAddr fallback without port",
			remoteAddr: "192.168.1.1",
			headers:    map[string]string{},
			expectedIP: "192.168.1.1",
		},
		{
			name:       "IPv6 with port",
			remoteAddr: "[::1]:12345",
			headers:    map[string]string{},
			expectedIP: "[::1]",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			req.RemoteAddr = tc.remoteAddr

			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}

			actualIP := getClientIP(req)
			assert.Equal(t, tc.expectedIP, actualIP)
		})
	}
}

func TestGetServiceFromEndpoint(t *testing.T) {
	testCases := []struct {
		endpoint        string
		expectedService string
	}{
		{"/warehouses", "warehouse-service"},
		{"/warehouses/123", "warehouse-service"},
		{"/warehouses/123/items", "warehouse-service"},
		{"/products", "product-service"},
		{"/products/456", "product-service"},
		{"/sections", "section-service"},
		{"/sections/789", "section-service"},
		{"/unknown", "api-service"},
		{"/random/endpoint", "api-service"},
		{"", "api-service"},
		{"/", "api-service"},
		{"//", "api-service"},
		{"/api/v1/warehouses", "api-service"}, // First segment is 'api'
	}

	for _, tc := range testCases {
		t.Run(tc.endpoint, func(t *testing.T) {
			actualService := getServiceFromEndpoint(tc.endpoint)
			assert.Equal(t, tc.expectedService, actualService)
		})
	}
}

func TestLoggingMiddleware_WithUserAgent(t *testing.T) {
	mockLogger := new(MockLogger)
	middleware := LoggingMiddleware(mockLogger)

	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"

	mockLogger.On("LogRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, userAgent, mock.Anything, mock.Anything, mock.Anything).Return()
	mockLogger.On("LogResponse", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify user agent in context
		ctxUserAgent := r.Context().Value("user_agent")
		assert.Equal(t, userAgent, ctxUserAgent)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", userAgent)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	mockLogger.AssertExpectations(t)
}

func TestLoggingMiddleware_WithoutUserAgent(t *testing.T) {
	mockLogger := new(MockLogger)
	middleware := LoggingMiddleware(mockLogger)

	mockLogger.On("LogRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, "", mock.Anything, mock.Anything, mock.Anything).Return()
	mockLogger.On("LogResponse", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	mockLogger.AssertExpectations(t)
}

func TestLoggingMiddleware_IntegrationTest(t *testing.T) {
	mockLogger := new(MockLogger)
	middleware := LoggingMiddleware(mockLogger)

	var capturedRequestID string
	var capturedStatusCode int
	var capturedExecutionTime int

	// Capture the logged values
	mockLogger.On("LogRequest", mock.Anything, "warehouse-service", "POST", "/warehouses", mock.AnythingOfType("string"), mock.AnythingOfType("string"), (*int)(nil), mock.AnythingOfType("string")).
		Run(func(args mock.Arguments) {
			capturedRequestID = args.Get(7).(string)
		}).Return()

	mockLogger.On("LogResponse", mock.Anything, "warehouse-service", mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.Anything).
		Run(func(args mock.Arguments) {
			capturedStatusCode = args.Get(3).(int)
			capturedExecutionTime = args.Get(4).(int)
		}).Return()

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify all context values are set
		assert.NotNil(t, r.Context().Value("request_id"))
		assert.Equal(t, "/warehouses", r.Context().Value("endpoint"))
		assert.Equal(t, "POST", r.Context().Value("method"))
		assert.NotNil(t, r.Context().Value("user_agent"))
		assert.NotNil(t, r.Context().Value("ip_address"))

		time.Sleep(5 * time.Millisecond) // Simulate some work
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 1, "name": "Test Warehouse"}`))
	}))

	req := httptest.NewRequest("POST", "/warehouses", strings.NewReader(`{"name": "Test"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "test-client/1.0")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, `{"id": 1, "name": "Test Warehouse"}`, w.Body.String())

	// Verify Request ID header
	responseRequestID := w.Header().Get("X-Request-ID")
	assert.NotEmpty(t, responseRequestID)
	assert.Equal(t, capturedRequestID, responseRequestID)

	// Verify captured values
	assert.Equal(t, http.StatusCreated, capturedStatusCode)
	assert.GreaterOrEqual(t, capturedExecutionTime, 5) // At least 5ms

	mockLogger.AssertExpectations(t)
}

func TestLoggingMiddleware_PanicRecovery(t *testing.T) {
	mockLogger := new(MockLogger)
	middleware := LoggingMiddleware(mockLogger)

	mockLogger.On("LogRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	// LogResponse may or may not be called depending on when panic occurs
	mockLogger.On("LogResponse", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return().Maybe()

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Should not panic the test
	assert.Panics(t, func() {
		handler.ServeHTTP(w, req)
	})
}

// Benchmark tests
func BenchmarkLoggingMiddleware(b *testing.B) {
	mockLogger := new(MockLogger)
	middleware := LoggingMiddleware(mockLogger)

	// Allow unlimited calls for benchmark
	mockLogger.On("LogRequest", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	mockLogger.On("LogResponse", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}
}
