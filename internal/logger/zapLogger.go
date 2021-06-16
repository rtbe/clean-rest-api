package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger is an wrapper for zap logger
// which satisfies logger interface
// ---
// link: https://github.com/uber-go/zap
type ZapLogger struct {
	*zap.Logger
}

// NewZapLogger returns a new zap logger.
func NewZapLogger() *ZapLogger {
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	logger := zap.New(core)
	zapLogger := &ZapLogger{
		logger.Named("REST API"),
	}
	return zapLogger
}

// Bunch of settings for zap logger.
func getEncoder() zapcore.Encoder {
	// Configure encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.CallerKey = "caller"
	// By default zap uses JSON encoder so let`s change it
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// Log info according to passed level and message.
func (z *ZapLogger) Log(level, message string) {
	switch strings.ToLower(level) {
	case "info":
		z.Info(message)
	case "error":
		z.Error(message)
	case "fatal":
		z.Fatal(message)
	}
}
