package util

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// TestLogger is a dummy logger for unit tests
type TestLogger struct {
	t testing.T
	mock.Mock
}

// Debug is a dummy logger method
func (l *TestLogger) Debug(args ...interface{}) { l.Called(args) }

// Info is a dummy logger method
func (l *TestLogger) Info(args ...interface{}) { l.Called(args) }

// Warn is a dummy logger method
func (l *TestLogger) Warn(args ...interface{}) { l.Called(args) }

// Error is a dummy logger method
func (l *TestLogger) Error(args ...interface{}) { l.Called(args) }

// Fatal is a dummy logger method
func (l *TestLogger) Fatal(args ...interface{}) { l.Called(args) }

// Debugf is a dummy logger method
func (l *TestLogger) Debugf(format string, args ...interface{}) { l.Called(args) }

// Infof is a dummy logger method
func (l *TestLogger) Infof(format string, args ...interface{}) { l.Called(args) }

// Warnf is a dummy logger method
func (l *TestLogger) Warnf(format string, args ...interface{}) { l.Called(args) }

// Errorf is a dummy logger method
func (l *TestLogger) Errorf(format string, args ...interface{}) { l.Called(args) }

// Fatalf is a dummy logger method
func (l *TestLogger) Fatalf(format string, args ...interface{}) { l.Called(args) }
