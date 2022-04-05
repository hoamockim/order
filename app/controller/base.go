package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"order/app/service"
	"order/pkg/errors"
	"os"
)

type Meta struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type BaseResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

func success(c *gin.Context, data interface{}) {
	requestId := c.GetHeader("X-Request-ID")
	if data != nil {
		jsonRes, _ := json.Marshal(data)
		fmt.Printf("requestId: %s, response: %s", requestId, jsonRes)
	}
	c.JSON(http.StatusOK, data)
}

func handleErr(c *gin.Context, errMeta *errors.ErrorMeta, serviceName string) {

}

func logDebug(data string) {
	_, _ = fmt.Fprintf(os.Stdout, "[GIN-debug] [Request Info] "+data)
}

var sysCtrl *sysController
var orderCtrl *orderController

func init() {
	sysCtrl = &sysController{}
	orderCtrl = &orderController{}
	orderCtrl.service = service.NewOrderService()
}
