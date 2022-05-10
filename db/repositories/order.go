package repositories

import (
	"fmt"
	"order/db"
	"order/db/models"
	"order/pkg/cache"
)

const cacheOrdersKey = "current_orders"

type QueryOrderRepository interface {
	GetOrderInfo(id string) (*models.OrderInfo, error)
	GetOrderList(customerId string) ([]*models.OrderInfo, error)
}

type CommandOrderRepository interface {
	CreateOrder(*models.OrderInfo) error
	CreateMultiOrder([]*models.OrderInfo) error
	UpdateOrderItems(orderCode string, items []*models.OrderItem) error
}

func (repo *defaultRepository) GetOrderInfo(id string) (*models.OrderInfo, error) {
	var orderInfo models.OrderInfo
	if cache.Instance().Get(fmt.Sprintf("order:%s", id), &orderInfo); len(orderInfo.OrderItems) > 0 {
		return &orderInfo, nil
	}
	return &orderInfo, nil
}

func (repo *defaultRepository) GetOrderList() ([]*models.OrderInfo, error) {
	var orderInfos []*models.OrderInfo
	if cache.Instance().Get(cacheOrdersKey, &orderInfos); len(orderInfos) > 0 {
		return orderInfos, nil
	}
	documents, err := db.GetAll(&models.OrderInfo{})
	if err != nil {
		return nil, err
	}
	orderInfos = makeOrderList(documents)
	cache.Instance().Set(cacheOrdersKey, &orderInfos, cache.DefaultHourTime)
	return orderInfos, nil
}

func makeOrderList(documents []interface{}) []*models.OrderInfo {
	var orderInfos []*models.OrderInfo
	for _, value := range documents {
		if doc, ok := value.(*models.OrderInfo); ok {
			orderInfos = append(orderInfos, doc)
		}
	}
	return orderInfos
}

func (repo *defaultRepository) CreateOrder(orderInfo *models.OrderInfo) (err error) {
	err = db.Save(orderInfo)
	return
}

func (repo *defaultRepository) CreateMultiOrder(orderInfoList []*models.OrderInfo) error {
	var documents []interface{}
	for order, _ := range orderInfoList {
		documents = append(documents, order)
	}
	return db.InsertBatch(orderInfoList[0].CollectionName(), documents)
}

func (repo *defaultRepository) UpdateOrderItems(orderCode string, items []*models.OrderItem) error {
	/*
		order := &models.OrderInfo{OrderItems: items}
			order.Id, _ = primitive.ObjectIDFromHex(orderCode)
			conditions := []string{"id"}
			updatedFields := []string{"items"}
			if err := db.Update(order, conditions, updatedFields); err != nil {
				return err
			}
	*/
	return nil
}
