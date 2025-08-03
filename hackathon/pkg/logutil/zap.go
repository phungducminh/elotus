package logutil

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	JSONLogFormat    = "json"
	ConsoleLogFormat = "console"
	//revive:disable:var-naming
	// Deprecated: Please use JSONLogFormat.
	JsonLogFormat = JSONLogFormat
	//revive:enable:var-naming
)

var DefaultLogFormat = JSONLogFormat


// CreateDefaultZapLogger: implementation taken from etcd: https://github.com/etcd-io/etcd/
func CreateDefaultZapLogger(level zapcore.Level) (*zap.Logger, error) {
	lcfg := DefaultZapLoggerConfig
	lcfg.Level = zap.NewAtomicLevelAt(level)
	c, err := lcfg.Build()
	if err != nil {
		return nil, err
	}
	return c, nil
}

var DefaultZapLoggerConfig = zap.Config{
	Level: zap.NewAtomicLevelAt(ConvertToZapLevel(DefaultLogLevel)),

	Development: false,
	Sampling: &zap.SamplingConfig{
		Initial:    100,
		Thereafter: 100,
	},

	Encoding: DefaultLogFormat,

	// copied from "zap.NewProductionEncoderConfig" with some updates
	EncoderConfig: zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,

		// Custom EncodeTime function to ensure we match format and precision of historic capnslog timestamps
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02T15:04:05.999999Z0700"))
		},

		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	},

	// Use "/dev/null" to discard all
	OutputPaths:      []string{"stderr"},
	ErrorOutputPaths: []string{"stderr"},
}

var DefaultLogLevel = "info"

// ConvertToZapLevel converts log level string to zapcore.Level.
func ConvertToZapLevel(lvl string) zapcore.Level {
	var level zapcore.Level
	if err := level.Set(lvl); err != nil {
		panic(err)
	}
	return level
}
