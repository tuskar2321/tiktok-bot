package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	cfg    zap.Config
	Logger *zap.SugaredLogger
)

func InitLogger() error {
	cfg = zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",
			LevelKey:   "level",
			TimeKey:    "ts",
			EncodeTime: zapcore.ISO8601TimeEncoder,
		},
	}

	l, err := cfg.Build()
	if err != nil {
		return err
	}
	defer l.Sync()
	Logger = l.Sugar()
	return nil
}
