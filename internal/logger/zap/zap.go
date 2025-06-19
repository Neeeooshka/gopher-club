package zap

import (
	"github.com/Neeeooshka/gopher-club/internal/logger"
	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
}

var Log = &ZapLogger{logger: zap.NewNop()}

func (l *ZapLogger) init(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	l.logger = zl
	return nil
}

func (l *ZapLogger) Log(rq logger.RequestData, rs logger.ResponseData) {
	l.logger.Info("receive new request",
		zap.String("URI", rq.URI),
		zap.String("method", rq.Method),
		zap.Duration("duration", rq.Duration),
		zap.Int("status", rs.Status),
		zap.Int("size", rs.Size),
	)
}

func (l *ZapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *ZapLogger) String(key string, val string) zap.Field {
	return zap.String(key, val)
}

func (l *ZapLogger) Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

func (l *ZapLogger) Error(err error) zap.Field {
	return zap.Error(err)
}

func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func NewZapLogger(level string) (*ZapLogger, error) {
	return Log, Log.init(level)
}
