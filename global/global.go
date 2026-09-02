package global

import (
	"database/sql"
	"golf-booking-go/pkg/logger"
	"golf-booking-go/pkg/setting"

	"github.com/redis/go-redis/v9"
)

var (
	Config *setting.Config
	Logger *logger.LoggerZap
	RDB    *redis.Client
	DB     *sql.DB
)
