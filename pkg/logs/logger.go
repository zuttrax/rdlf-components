package logs

import (
	"fmt"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	ErrorLevel = zapcore.ErrorLevel
)

type Logger interface {
	Info(string)
	Error(error)
	Fatal(error)
}

type Log struct {
	logger *zap.SugaredLogger
}

func InitializeLog(appName string, level string) Log {
	writer := createLogWriterByAppName(appName)

	encoder := getLogEncoder()

	zapLevel := getLogLevel(level)

	core := zapcore.NewCore(encoder, writer, zapLevel)

	logger := zap.New(core).Sugar()

	return Log{
		logger: logger,
	}
}

func createLogWriterByAppName(appName string) zapcore.WriteSyncer {
	fileName := fmt.Sprintf("/rdlf/%s.log", appName)

	lumberJackWriter := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    10,
		MaxBackups: 2,
		MaxAge:     5,
		Compress:   false,
	}

	return zapcore.AddSync(lumberJackWriter)
}

func getLogEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func (l Log) Info(msg string) {
	l.logger.Info(msg)
}

func (l Log) Error(err error) {
	l.logger.Error(err.Error())
}

func (l Log) Fatal(err error) {
	l.logger.Fatal(err.Error())
}
