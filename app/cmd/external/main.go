package main

import (
	"github.com/gin-gonic/gin"
	"order/app"
	"order/app/controller"
	"order/pkg/configs"
	"os"
)

var r *gin.Engine

func main() {
	ordV1 := r.Group("tipee/v1/")
	{
		commonService := ordV1.Group("sys/")
		{
			commonService.GET("health-check", controller.HealthCheck)
		}
		orderService := ordV1.Group("order/")
		{
			orderService.GET("detail", controller.GetOrder)
			orderService.POST("", controller.CreateOrder)
		}
	}

	if err := r.Run(configs.AppURL()); err != nil {
		os.Exit(1)
	}
}

func init() {
	r = gin.New()
	app.InitDomain()
}
