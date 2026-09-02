package initialize

import (
	"golf-booking-go/global"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	// r := routers.NewRouter()
	var r *gin.Engine
	if global.Config.Server.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}
	// manageRouter := routers.RouterGroupApp.Manage
	// userRouter := routers.RouterGroupApp.User

	// MainGroup := r.Group("/v1")
	// {
	// 	MainGroup.GET("/health")
	// }
	// {
	// 	manageRouter.InitAdminRouter(MainGroup)
	// 	manageRouter.InitUserRouter(MainGroup)
	// }
	// {
	// 	userRouter.InitUserRouter(MainGroup)
	// 	userRouter.InitProductRouter(MainGroup)
	// }

	return r
}
