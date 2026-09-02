package global

import (
	"golf-booking-go/pkg/logger"
	"golf-booking-go/pkg/setting"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Config *setting.Config
	Logger *logger.LoggerZap
	RDB    *redis.Client
	DB     *gorm.DB
)
