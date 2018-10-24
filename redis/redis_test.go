package redis

import (
	"log"
	"os"
	"testing"
)

var client *Cache

type mlogger struct {
	*log.Logger
}

func (l *mlogger) Debug(format string, args ...interface{}) {
	l.Printf(format, args...)
}

func (l *mlogger) Info(format string, args ...interface{}) {
	l.Printf(format, args...)
}

func (l *mlogger) Error(format string, args ...interface{}) {
	l.Printf(format, args...)
}

func init() {
	cache, err := NewCache(Options{
		NodeAddr:[]string{"127.0.0.1:6379"},
		Db:1,
		Password:"",
		Logger:&mlogger{log.New(os.Stdout, "", log.LstdFlags)},
	})
	if err != nil {
		panic(err)
	}
	client = cache
}


func TestCache_PatternDel(t *testing.T) {
	client.Set("t", "test")
	client.Set("t2", "test2")
	a := ""
	client.Get("t", &a)
	log.Printf("t="+a)
	err := client.PatternDel("t*")
	if err != nil {
		t.Error(err)
	}
}
