package initialize

import (
	"fmt"
	"golf-booking-go/global"
	"golf-booking-go/internal/router"

	"go.uber.org/zap"
)

func Run() {
	LoadConfig()
	fmt.Println("Loaded config: ", global.Config)
	InitLogger()
	global.Logger.Info("Config log ok!!", zap.String("ok", "success"))
	InitMysql()
	InitRedis()
	InitMiddlewares()

	r := InitRouter()

	apiV1Group := r.Group("/api/v1")

	router.InitRouter(apiV1Group)

	r.Run(":8002")
}
