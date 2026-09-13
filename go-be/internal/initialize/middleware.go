package initialize

import "golf-booking-go/internal/middlewares"

func InitMiddlewares() {
	middlewares.InitAuthenticationMiddleware()
}
