package controller

import (
	"github.com/gin-gonic/gin"
	"order/app/dto"
	"order/app/service"
	"order/pkg/errors"
)

type orderController struct {
	service interface {
		service.CommandOrderService
		service.QueryOrderService
	}
}

func CreateOrder(ctx *gin.Context) {
	var req dto.OrderRequest
	if err := ctx.Bind(&req); err != nil {

	}
	err := orderCtrl.service.CreateOrder(&req)
	if err != nil {
		handleErr(ctx, &errors.ErrorMeta{HttpCode: 400}, "create-order")
	}
	success(ctx, &BaseResponse{Meta{Code: "201", Message: "create order successful"}, nil})
}

func GetOrder(ctx *gin.Context) {
	code := ctx.Param("code")
	if code == "" {

	}
	orderData, err := orderCtrl.service.GetOrderDetail(code)
	if err != nil {

	}
	success(ctx, &BaseResponse{
		Meta: Meta{"200", ""},
		Data: &orderData,
	})
}
