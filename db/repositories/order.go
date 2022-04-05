package repositories

import "order/db/models"

type QueryOrderRepository interface {
	GetOrderInfo() (*models.OrderInfo, error)
}

type CommandOrderRepository interface {
	CreateOrder(*models.OrderInfo) error
}

func (repo *defaultRepository) CreateOrder(orderInfo *models.OrderInfo) ( err error) {
	err = repo.save(orderInfo)
	return
}
