package global

import (
	"github.com/redis/go-redis/v9"
	"github.com/xxandjg/ginbase/entity"
	"go.uber.org/zap"
	"strconv"
)

var (
	RedisDB = new(redis.Client)
)

func InitRedis(config *entity.RedisConfig) Error {

	options := &redis.Options{
		Addr:     config.Host + ":" + strconv.Itoa(config.Port),
		Network:  config.Network,
		Password: config.Password,
		DB:       config.DB,
	}
	RedisDB = redis.NewClient(options)
	zap.L().Sugar().Info("init RedisDB success")
	return SUCCESS
}
