package interfaces

// Logger is an interface to viper
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

// TestLogger is a dummy logger for unit tests
type TestLogger struct{}

// Debug is a dummy logger method
func (l *TestLogger) Debug(args ...interface{}) {}

// Info is a dummy logger method
func (l *TestLogger) Info(args ...interface{}) {}

// Warn is a dummy logger method
func (l *TestLogger) Warn(args ...interface{}) {}

// Error is a dummy logger method
func (l *TestLogger) Error(args ...interface{}) {}

// Fatal is a dummy logger method
func (l *TestLogger) Fatal(args ...interface{}) {}

// Debugf is a dummy logger method
func (l *TestLogger) Debugf(format string, args ...interface{}) {}

// Infof is a dummy logger method
func (l *TestLogger) Infof(format string, args ...interface{}) {}

// Warnf is a dummy logger method
func (l *TestLogger) Warnf(format string, args ...interface{}) {}

// Errorf is a dummy logger method
func (l *TestLogger) Errorf(format string, args ...interface{}) {}

// Fatalf is a dummy logger method
func (l *TestLogger) Fatalf(format string, args ...interface{}) {}
