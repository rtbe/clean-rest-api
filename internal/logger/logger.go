// Package logger provides wrapper interface for logging to make it implementation-independent.
// As well as implementation of it with zap logger.
package logger

// Logger is an interface for loggers
type Logger interface {
	Log(level, message string)
}
