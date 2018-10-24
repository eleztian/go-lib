package redis

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
}

type Cache struct {
	client *redis.Client
	ttl    time.Duration

	logger *mlog
}

type Options struct {
	NodeAddr       []string
	SentinelEnable bool
	MasterName     string
	PoolSize       int
	Db             int
	Password       string
	Logger         Logger

	TTL          time.Duration
	DialTimeOut  time.Duration
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	PoolTimeout  time.Duration
}

func NewCache(op Options) (*Cache, error) {
	var r = &Cache{
		ttl: op.TTL,
	}
	r.logger = &mlog{op.Logger}
	if op.SentinelEnable {
		r.client = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    op.MasterName,
			SentinelAddrs: op.NodeAddr,
			Password:      op.Password,
			DialTimeout:   op.DialTimeOut,
			ReadTimeout:   op.ReadTimeout,
			WriteTimeout:  op.WriteTimeout,
			PoolSize:      op.PoolSize,
			PoolTimeout:   op.PoolTimeout,
		})
	} else {
		r.client = redis.NewClient(&redis.Options{
			Addr:         op.NodeAddr[0],
			Password:     op.Password,
			DialTimeout:  op.DialTimeOut,
			ReadTimeout:  op.ReadTimeout,
			WriteTimeout: op.WriteTimeout,
			PoolSize:     op.PoolSize,
			PoolTimeout:  op.PoolTimeout,
			DB:           op.Db,
		})
	}

	_, err := r.client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Cache) Set(key string, value interface{}) error {
	return r.SetWithExp(key, value, 0)
}

func (r *Cache) SetWithExp(key string, value interface{}, exp time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		r.logger.Error("marshal[%s = %#v] error:%v", key, value, err)
		return err
	}
	err = r.client.Set(key, string(b), exp).Err()
	if err != nil {
		r.logger.Error("set [%s = %#v] error:%v", key, value, err)
		return err
	}
	r.logger.Debug("set [%s = %v] success", key, value)
	return nil
}

// Get value must a pointer
func (r *Cache) Get(key string, value interface{}) (err error) {
	sc := r.client.Get(key)
	if err = sc.Err(); err != nil {
		if err == redis.Nil {
			r.logger.Debug("miss %s", key)
			value = nil
			return nil
		} else {
			r.logger.Error("get [%s=%#v] error:%v", key, value, err)
			return err
		}
	}
	a := new(string)
	if err = sc.Scan(a); err != nil {
		r.logger.Error("scan [%s=%#v] error:%v", key, value, err)
		return err
	}
	err = json.Unmarshal([]byte(*a), value)
	if err != nil {
		r.logger.Error("unmarshal [%s=%#v] error:%v", key, value, err)
		return err
	}
	r.logger.Debug("hit %s", key)
	return nil
}

func (r *Cache) Del(key string) error {
	err := r.client.Del(key).Err()
	if err != nil && err != redis.Nil {
		r.logger.Error("del [%s=%#v] error:%v", key, err)
		return err
	}
	r.logger.Debug("del %s", key)
	return nil
}

const delKeyByPattern = `
local ks = redis.call("keys",KEYS[1])
if table.getn(ks) ~= 0 then 
	redis.call("del", unpack(ks))
end
`

func (r *Cache) PatternDel(pattern string) error {
	_, err := r.Lua(delKeyByPattern, pattern)
	return err
}

func (r *Cache) Lua(lua string, keys ...string) (interface{}, error) {
	c := redis.NewScript(lua).Run(r.client, keys)
	if err := c.Err(); err != nil && err != redis.Nil {
		r.logger.Error("lua [%s] error:%v", lua, err)
		return "", err
	}
	r.logger.Debug("lua [%s]", lua)
	rsp, _ := c.Result()
	return rsp, nil
}

func (r *Cache) Close() error {
	return r.client.Close()
}


type mlog struct {
	log Logger
}

func (m *mlog) Debug(format string, args ...interface{}) {
	if m.log != nil {
		m.log.Debug(format, args...)
	}
}

func (m *mlog) Info(format string, args ...interface{}) {
	if m.log != nil {
		m.log.Info(format, args...)
	}
}

func (m *mlog) Error(format string, args ...interface{}) {
	if m.log != nil {
		m.log.Error(format, args...)
	}
}
