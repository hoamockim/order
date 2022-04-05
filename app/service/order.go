package service

import (
	"order/app/dto"
	"order/db/models"
)

type QueryOrderService interface {
	GetOrderDetail(code string) (*dto.OrderData, error)
}

type CommandOrderService interface {
	CreateOrder(req *dto.OrderRequest) error
}

func (srv *defaultService) CreateOrder(req *dto.OrderRequest) error {
	//main follow
	orderInfo := models.OrderInfo{
		Address:    req.AddressCode,
		CustomerId: req.CustomerId,
	}
	//store event into event-sourcing
	err := srv.OrderRepo.CreateOrder(&orderInfo)
	if err != nil {
		return err
	}
	//push event create-order
	return nil
}

func (srv *defaultService) GetOrderDetail(code string) (oderData *dto.OrderData, err error) {
	return
}
