package dao

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/ningchengzeng/goim/internal/logic/conf"

	"github.com/nsqio/go-nsq"
)

// Dao dao.
type Dao struct {
	c           *conf.Config
	nsqPub      *nsq.Producer
	redis       *redis.Pool
	redisExpire int32
}

// New new a dao and return.
func New(c *conf.Config) *Dao {
	d := &Dao{
		c:           c,
		nsqPub:      newNsqPub(c.Nsq),
		redis:       newRedis(c.Redis),
		redisExpire: int32(time.Duration(c.Redis.Expire) / time.Second),
	}
	return d
}

func newNsqPub(c *conf.Nsq) *nsq.Producer {
	cfg := nsq.NewConfig()

	producer, err := nsq.NewProducer(c.Address, cfg)
	if err != nil {
		return nil
	}
	return producer
}

func newRedis(c *conf.Redis) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     c.Idle,
		MaxActive:   c.Active,
		IdleTimeout: time.Duration(c.IdleTimeout),
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(c.Network, c.Addr,
				redis.DialConnectTimeout(time.Duration(c.DialTimeout)),
				redis.DialReadTimeout(time.Duration(c.ReadTimeout)),
				redis.DialWriteTimeout(time.Duration(c.WriteTimeout)),
				redis.DialPassword(c.Auth),
			)
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}

// Close close the resource.
func (d *Dao) Close() error {
	return d.redis.Close()
}

// Ping dao ping.
func (d *Dao) Ping(c context.Context) error {
	return d.pingRedis(c)
}
