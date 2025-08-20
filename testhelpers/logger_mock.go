package testhelpers

import (
	"context"

	pkglogger "github.com/varobledo_meli/W17-G10-Bootcamp.git/pkg/logger"
)

// Ensure TestLogger implements the logger.Logger interface
var _ pkglogger.Logger = (*TestLogger)(nil)

// TestLogger is a no-op implementation of the logger.Logger interface for tests.
type TestLogger struct{}

func NewTestLogger() *TestLogger { return &TestLogger{} }

func (l *TestLogger) Debug(ctx context.Context, service, message string, metadata ...map[string]interface{})                 {}
func (l *TestLogger) Info(ctx context.Context, service, message string, metadata ...map[string]interface{})                  {}
func (l *TestLogger) Warning(ctx context.Context, service, message string, metadata ...map[string]interface{})               {}
func (l *TestLogger) Error(ctx context.Context, service, message string, err error, metadata ...map[string]interface{})      {}
func (l *TestLogger) Fatal(ctx context.Context, service, message string, err error, metadata ...map[string]interface{})       {}
func (l *TestLogger) LogRequest(ctx context.Context, service, method, endpoint, userAgent, ipAddress string, userID *int, requestID string) {}
func (l *TestLogger) LogResponse(ctx context.Context, service, requestID string, statusCode, executionTimeMs int, metadata ...map[string]interface{}) {}
