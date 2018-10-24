package vlog2

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"../context"
)

type Logger struct {
	*zap.Logger
	zap.AtomicLevel
	callSkip int
}

func (l *Logger) GetLevel() int8 {
	return int8(l.AtomicLevel.Level()) + 1
}

func (l *Logger) SetLevel(level int8) {
	l.AtomicLevel.SetLevel(zapcore.Level(level - 1))
}

func (l *Logger) output(callSkip int, level zapcore.Level, msg string, fields ...zapcore.Field) {
	callerSkip := l.callSkip + 1
	nl := l.WithOptions(
		zap.AddCaller(),
		zap.AddCallerSkip(int(callerSkip)),
	).With(zap.Int64("timeStamp", time.Now().UTC().UnixNano()))
	switch level {
	case zapcore.DebugLevel:
		nl.Debug(msg, fields...)
	case zapcore.InfoLevel:
		nl.Info(msg, fields...)
	case zapcore.WarnLevel:
		nl.Warn(msg, fields...)
	case zapcore.ErrorLevel:
		nl.Error(msg, fields...)
	case zapcore.DPanicLevel:
		nl.DPanic(msg, fields...)
	case zapcore.PanicLevel:
		nl.Panic(msg, fields...)
	case zapcore.FatalLevel:
		nl.Fatal(msg, fields...)
	default:
	}
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...zapcore.Field) {
	if ctx != nil {
		fields = append(fields, zap.String("traceID", ctx.String()))
	}
	l.output(l.callSkip+1, zapcore.DebugLevel, msg, fields...)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zapcore.Field) {
	if ctx != nil {
		fields = append(fields, zap.String("traceID", ctx.String()))
	}
	l.output(l.callSkip+1, zapcore.InfoLevel, msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...zapcore.Field) {
	if ctx != nil {
		fields = append(fields, zap.String("traceID", ctx.String()))
	}
	l.output(l.callSkip+1, zapcore.WarnLevel, msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, err error, fields ...zapcore.Field) {
	if ctx != nil {
		fields = append(fields, zap.String("traceID", ctx.String()), zap.Error(err))
	}
	l.output(l.callSkip+1, zapcore.ErrorLevel, msg, fields...)
}

func (l *Logger) Panic(ctx context.Context, msg string, fields ...zapcore.Field) {
	if ctx != nil {
		fields = append(fields, zap.String("traceID", ctx.String()))
	}
	l.output(l.callSkip+1, zapcore.PanicLevel, msg, fields...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zapcore.Field) {
	if ctx != nil {
		fields = append(fields, zap.String("traceID", ctx.String()))
	}
	l.output(l.callSkip+1, zapcore.FatalLevel, msg, fields...)
}

func NewLogger(writer zapcore.WriteSyncer, name string,
	level int8, callSkip int, m map[string]interface{}) (*Logger, error) {

	zl := zap.NewAtomicLevel()
	zl.SetLevel(zapcore.Level(level - 1))

	fs := make([]zap.Field, 0, len(m))
	for k, v := range m {
		fs = append(fs, zap.Any(k, v))
	}
	log := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "S",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "stacktraceKey",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     utcTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		}), writer, zl)).Named(name).With(fs...)

	return &Logger{log, zl, callSkip}, nil
}

func utcTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format("2006-01-02T15:04:05.00000"))
}
