package zap

import (
	"github.com/zedisdog/ty/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLog() log.ILog {
	var zapConfig = zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	z, err := zapConfig.Build()
	z = z.WithOptions(zap.AddCallerSkip(1))
	if err != nil {
		panic(log.Wrap(err, "new zap failed"))
	}

	return &Log{
		zap: z,
	}
}

type Log struct {
	zap *zap.Logger
}

func (l Log) convertFields(fields []*log.Field) (zapFields []zap.Field) {
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Name, field.Value))
	}
	return
}

func (l Log) Trace(msg string, fields ...*log.Field) {
	l.zap.Debug("trace: "+msg, l.convertFields(fields)...)
}

func (l Log) Debug(msg string, fields ...*log.Field) {
	l.zap.Debug(msg, l.convertFields(fields)...)
}

func (l Log) Info(msg string, fields ...*log.Field) {
	l.zap.Info(msg, l.convertFields(fields)...)
}

func (l Log) Warn(msg string, fields ...*log.Field) {
	l.zap.Warn(msg, l.convertFields(fields)...)
}

func (l Log) Error(msg string, fields ...*log.Field) {
	l.zap.Error(msg, l.convertFields(fields)...)
}

func (l Log) Fatal(msg string, fields ...*log.Field) {
	l.zap.Fatal(msg, l.convertFields(fields)...)
}
