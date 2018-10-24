package logger

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"testing"
)

type RedID struct {
	s string
}

func (r RedID) String() string {
	return r.s
}

func TestReleaseLog(t *testing.T) {
	logger := New(Options{
		Name:      "test",
		Outer:     os.Stdout,
		Mode:      LOG_MODE_RELEASE,
		CallSkip:  1,
		CommonMap: map[string]interface{}{"H": "127.0.0.1"},
		Level:     Lerror,
	})
	logger.Debug(RedID{"123"}, "this is a test", zap.String("t", "t"), zap.Int("aaa", 123))
	logger.Error(RedID{"123"}, "this is a test", fmt.Errorf("error"))
	logger.Info(RedID{"123"}, "this is a test")
	logger.Warn(RedID{"123"}, "this is a test")
	logger.Panic(RedID{"123"}, "this is a test")
}
