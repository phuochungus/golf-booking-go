package initialize

import (
	"golf-booking-go/global"
	"golf-booking-go/pkg/logger"
)

func InitLogger() {
	global.Logger = logger.NewLogger(global.Config.Logger)
}
