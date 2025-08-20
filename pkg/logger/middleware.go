package logger

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)



// LoggingMiddleware creates a complete middleware for HTTP request logging
func LoggingMiddleware(logger Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Generate a unique request ID using UUID v4
			requestID := uuid.New().String() // Extract request information
			method := r.Method
			endpoint := r.URL.Path
			service := getServiceFromEndpoint(endpoint)
			userAgent := r.UserAgent()
			ipAddress := getClientIP(r)

			// TODO: Extract user_id from authentication context when implemented
			var userID *int = nil

			// Add all information to context so it's available in logs
			ctx := context.WithValue(r.Context(), "request_id", requestID)
			ctx = context.WithValue(ctx, "endpoint", endpoint)
			ctx = context.WithValue(ctx, "method", method)
			ctx = context.WithValue(ctx, "user_agent", userAgent)
			ctx = context.WithValue(ctx, "ip_address", ipAddress)
			r = r.WithContext(ctx)

			// Log request start using the new LogRequest method
			logger.LogRequest(ctx, service, method, endpoint, userAgent, ipAddress, userID, requestID)

			// Create wrapper to capture status code
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Add Request ID to response header
			ww.Header().Set("X-Request-ID", requestID)

			// Execute the handler
			next.ServeHTTP(ww, r)

			// Calculate execution time
			duration := time.Since(start)
			statusCode := ww.Status()
			executionTimeMs := int(duration.Milliseconds())

			// Add response info to context
			ctx = context.WithValue(ctx, "status_code", statusCode)
			ctx = context.WithValue(ctx, "execution_time_ms", executionTimeMs)

			// Log request end using the new LogResponse method
			metadata := map[string]interface{}{
				"bytes_written": ww.BytesWritten(),
			}
			logger.LogResponse(ctx, service, requestID, statusCode, executionTimeMs, metadata)
		})
	}
}

// getClientIP extracts the real client IP considering proxies
func getClientIP(r *http.Request) string {
	// Check proxy headers first
	if xForwardedFor := r.Header.Get("X-Forwarded-For"); xForwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs, we take the first one
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	if xRealIP := r.Header.Get("X-Real-IP"); xRealIP != "" {
		return xRealIP
	}

	// Fallback to RemoteAddr
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

// getServiceFromEndpoint determines the service based on the endpoint
func getServiceFromEndpoint(endpoint string) string {
	parts := strings.Split(strings.Trim(endpoint, "/"), "/")
	if len(parts) == 0 {
		return "api"
	}

	// Basic mapping of endpoints to services
	switch parts[0] {
	case "warehouses":
		return "warehouse-service"
	case "products":
		return "product-service"
	case "sections":
		return "section-service"
	default:
		return "api-service"
	}
}
