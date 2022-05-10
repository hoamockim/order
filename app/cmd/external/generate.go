package main

import "github.com/gin-gonic/gin"

// Create produce a new order
// @Summary Produce an order
// @Schemes http
// @Description create an order
// @Tags tipee
// @Accept json
// @Produce json
// @Success 200 {struct} dto.BaseResponse
// @Param dto.OrderRequest required true
// @Router /order [post]
func Create(ctx *gin.Context) {}

// GetDetail get information of the order
// @Summary Get information of the order
// @Schemes http
// @Description Get information of the order
// @Tags tipee
// @Accept json
// @Produce json
// @Success 200 {struct} dto.BaseResponse
// @Param dto.OrderRequest required true
// @Router /order [post]
func GetDetail(ctx *gin.Context) {

}
