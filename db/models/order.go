package models

import (
	"go.mongodb.org/mongo-driver/bson"
)

type OrderInfo struct {
	OrderCode  string      `json:"order_code"`
	CustomerId string         `json:"customer_id"`
	Address    string      `json:"address"`
	OrderItems []OrderItem `json:"order_items"`
}

type OrderItem struct {
	ItemCode     string `json:"item_code"`
	ShopCode     string `json:"shop_code"`
	MerchantCode string `json:"merchant_code"`
	Number       int    `json:"number"`
}

func (orderInfo *OrderInfo) GetDocumentName() string {
	return "orders"
}

func (orderInfo *OrderInfo) IsCached() bool {
	return true
}

func (orderInfo *OrderInfo) GenerateBson() bson.D {
	//items, _ := json.Marshal(orderInfo.OrderItems)
	return bson.D{{"code", orderInfo.OrderCode}, {"items", orderInfo.OrderItems}}
}

func (orderInfo *OrderInfo) Validate() bool {
	return true
}
