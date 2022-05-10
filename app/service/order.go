package service

import (
	"order/app/dto"
	"order/db/models"
)

type QueryOrderService interface {
	GetDetail(code string) (*dto.OrderData, error)
}

type CommandOrderService interface {
	Create(req *dto.OrderRequest) error
	Cancel(req *dto.OrderProcessRequest) error
}

func (srv *defaultService) Create(req *dto.OrderRequest) error {
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

func (srv *defaultService) Cancel(req *dto.OrderProcessRequest) error {
	return nil
}

func (srv *defaultService) GetDetail(code string) (oderData *dto.OrderData, err error) {
	return
}
