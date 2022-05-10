package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swf "github.com/swaggo/files"
	gsw "github.com/swaggo/gin-swagger"
	"order/app"
	"order/app/cmd/external/docs"
	"order/app/controller"
	"order/app/middleware"
	"order/pkg/configs"
	"os"
	"strings"
)

const (
	AppExitCode = 1
)

var r *gin.Engine
var auth gin.HandlerFunc

func init() {
	r = gin.New()
	app.InitDomain()
	auth = middleware.Authorize()
	docs.SwaggerInfo.BasePath = "/tipee/order/v1/"
}

func main() {
	// public api interface
	ordV1 := r.Group("tipee/v1/")
	{
		sysService := ordV1.Group("sys/")
		{
			sysService.GET("health-check", controller.HealthCheck)
		}
		ordV1.Use(auth)
		orderService := ordV1.Group("order/")
		{
			orderService.GET("detail", controller.GetDetail)
			orderService.POST("", controller.Create)
		}
	}
	// swagger document
	swagger := fmt.Sprintf("http://%s/swagger/doc.json", configs.AppURL())
	r.GET("/swagger/*any", gsw.WrapHandler(swf.Handler, gsw.URL(swagger)))
	if err := r.Run(configs.AppURL()); err != nil {
		fatalf(err.Error())
	}
}

func fatalf(format string, values ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	_, _ = fmt.Fprintf(os.Stderr, "[GIN-debug] [ERROR] "+format, values...)
	os.Exit(AppExitCode)
}
