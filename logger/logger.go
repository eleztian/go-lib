package logger

import (
	"os"

	"go.uber.org/zap/zapcore"

	"./context"
	"./stdlog"
	"./vlog2"
)

const (
	Ldebug = 0
	Linfo  = 1
	Lwarn  = 2
	Lerror = 3
	Lpanic = 4
	Lfatal = 6
	Lnone  = 7
)

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...zapcore.Field)
	Info(ctx context.Context, msg string, fields ...zapcore.Field)
	Warn(ctx context.Context, msg string, fields ...zapcore.Field)
	Error(ctx context.Context, msg string, err error, fields ...zapcore.Field)
	Panic(ctx context.Context, msg string, fields ...zapcore.Field)
	Fatal(ctx context.Context, msg string, fields ...zapcore.Field)

	SetLevel(int8)
	GetLevel() int8
}

const (
	LOG_MODE_DEBUG   = "DEBUG"
	LOG_MODE_RELEASE = "RELEASE"
)

type Options struct {
	Name      string
	Mode      string
	Level     int8
	CallSkip  int
	CommonMap map[string]interface{}
	Outer     *os.File
}

func New(op Options) Logger {
	if op.Mode == "" || op.Mode == LOG_MODE_DEBUG {
		l, err := stdlog.NewLogger(op.Outer, op.Name, op.Level, op.CallSkip, op.CommonMap)
		if err != nil {
			panic(err)
		}
		return l
	} else if op.Mode == LOG_MODE_RELEASE {
		l, err := vlog2.NewLogger(op.Outer, op.Name, op.Level-1, op.CallSkip, op.CommonMap)
		if err != nil {
			panic(err)
		}
		return l
	}
	panic("unknown log mod")
}

var log Logger

func Init(op Options) {
	log = New(op)
}

func Debug(ctx context.Context, msg string, fields ...zapcore.Field) {
	log.Debug(ctx, msg, fields...)
}
func Info(ctx context.Context, msg string, fields ...zapcore.Field) {
	log.Info(ctx, msg, fields...)
}
func Warn(ctx context.Context, msg string, fields ...zapcore.Field) {
	log.Warn(ctx, msg, fields...)
}
func Error(ctx context.Context, msg string, err error, fields ...zapcore.Field) {
	log.Error(ctx, msg, err, fields...)
}
func Panic(ctx context.Context, msg string, fields ...zapcore.Field) {
	log.Panic(ctx, msg, fields...)
}
func Fatal(ctx context.Context, msg string, fields ...zapcore.Field) {
	log.Fatal(ctx, msg, fields...)
}
func SetLevel(l int8) {
	log.SetLevel(l)
}
func GetLevel() int8 {
	return log.GetLevel()
}
