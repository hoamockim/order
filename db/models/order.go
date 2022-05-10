package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderInfo struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	CustomerId string             `json:"customer_id" bson:"customer_id"`
	Address    string             `json:"address" bson:"address	"`
	OrderItems []*OrderItem       `json:"order_items" bson:"order_items"`
}

type OrderItem struct {
	Id           string `json:"id" bson:"id"`
	ShopCode     string `json:"shop_code" bson:"shop_code"`
	MerchantCode string `json:"merchant_code" bson:"merchant_code"`
	Number       int    `json:"number" bson:"number"`
}

func (orderInfo *OrderInfo) CollectionName() string {
	return "orders"
}

func (orderInfo *OrderInfo) IsCached() bool {
	return true
}

func (orderInfo *OrderInfo) Validate() bool {
	return true
}

func (orderInfo *OrderInfo) MakeCondition(filterFields []string) map[string]interface{} {
	conditions := make(map[string]interface{})
	for i := 0; i < len(filterFields); i++ {
		switch filterFields[i] {
		case "id":
			conditions["_id"] = orderInfo.Id
		case "customer_id":
			conditions["customer_id"] = orderInfo.CustomerId
		}
	}
	return conditions
}

func (orderInfo *OrderInfo) MakeUpdated(fields []string) []primitive.E {
	var updated []primitive.E
	for i := 0; i < len(fields); i++ {
		switch fields[i] {
		case "customer_id":
			updated = append(updated, primitive.E{Key: fields[i], Value: orderInfo.CustomerId})
		case "address":
			updated = append(updated, primitive.E{Key: fields[i], Value: orderInfo.Address})
		case "items":
			updated = append(updated, primitive.E{Key: fields[i], Value: orderInfo.OrderItems})
		}
	}
	return updated
}
