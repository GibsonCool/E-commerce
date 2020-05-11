package datasourse

import (
	"E-commerce/common/conf"
	"github.com/kataras/iris/sessions/sessiondb/redis"
	"github.com/kataras/iris/sessions/sessiondb/redis/service"
)

var redisInstance *redis.Database

func GetRedisInstance() *redis.Database {
	if redisInstance == nil {
		redisInstance = initRedis()
	}
	return redisInstance
}

func initRedis() *redis.Database {
	var dataBase *redis.Database

	rd := conf.RedisSetting
	dataBase = redis.New(service.Config{
		Network:     rd.NetWork,
		Addr:        rd.Host,
		Password:    rd.Password,
		Database:    "",
		MaxIdle:     rd.MaxIdle,
		MaxActive:   rd.MaxActive,
		IdleTimeout: service.DefaultRedisIdleTimeout,
		Prefix:      rd.Prefix,
	})

	return dataBase
}
