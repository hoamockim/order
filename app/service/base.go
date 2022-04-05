package service

import "order/db/repositories"

type defaultService struct {
	OrderRepo interface {
		repositories.CommandOrderRepository
	}
}

var dfSer orderService

type orderService interface {
	QueryOrderService
	CommandOrderService
}

func NewOrderService() orderService {
	return dfSer
}

func initService() interface {
	QueryOrderService
	CommandOrderService
} {
	return &defaultService{repositories.New()}
}

func init() {
	dfSer = initService()
}
