package models

import (
	"fmt"
	"time"

	"../modules/setting"
	"github.com/garyburd/redigo/redis"
	"github.com/weisd/log"
)

var (
	RedisPools map[string]*redis.Pool
)

func init() {
	RedisPools = make(map[string]*redis.Pool)
}

func Redis(name string) *redis.Pool {
	p, ok := RedisPools[name]
	if !ok {
		panic(fmt.Errorf("Unknown Redis name %s", name))
	}

	return p
}

func InitRedisPools() {
	for name, conf := range setting.Cfg.Redis {
		log.Debug("InitRedisPools name %s, conf %v", name, conf)
		RedisPools[name] = newRedis(conf)
	}
	log.Debug("初始化redis连接池 done \n %v", RedisPools)

}

func newRedis(conf setting.RedisConfig) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     conf.MAX_IDLE,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.ADDR)
			if err != nil {
				return nil, err
			}
			if len(conf.PASSWD) > 0 {
				if _, err := c.Do("AUTH", conf.PASSWD); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func RedisCheckConn() {
	for name, pool := range RedisPools {
		go func(name string, pool *redis.Pool) {
			log.Debug("check Redis conn name : %s", name)

			conn := pool.Get()
			res, err := conn.Do("PING")
			if err != nil {
				log.Error(4, "redis name %s, conn failed  %v", name, err)
			} else {
				log.Debug("redis name %s, conn ok %v", name, res)
			}
			conn.Close()
		}(name, pool)

	}
}
