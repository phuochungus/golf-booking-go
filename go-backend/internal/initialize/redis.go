package initialize

import (
	"context"
	"fmt"
	"golf-booking-go/global"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()

func InitRedis() {
	r := global.Config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", r.Host, r.Port),
		Password: r.Password,
		DB:       r.DB,
		PoolSize: r.PoolSize,
	})

	_, err := rdb.Ping(ctx).Result()

	if err != nil {
		global.Logger.Error("Redis connection error", zap.Error(err))
		return
	}

	global.RDB = rdb
	global.Logger.Info("Redis connection successful")
}
