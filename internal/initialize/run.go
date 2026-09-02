package initialize

import (
	"fmt"
	"golf-booking-go/global"

	"go.uber.org/zap"
)

func Run() {
	LoadConfig()
	fmt.Println("Loaded config: ", global.Config)
	InitLogger()
	global.Logger.Info("Config log ok!!", zap.String("ok", "success"))
	InitMysql()
	InitRedis()

	r := InitRouter()

	r.Run(":8002")
}
