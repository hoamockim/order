package dto

type OrderStatus int

type OrderRequest struct {
	Items       []OrderItem `json:"items"`
	CustomerId  string      `json:"customer_id"`
	AddressCode string      `json:"address_code"`
	CreatedAt   int         `json:"created_at"`
}

type OrderProcessRequest struct {
	Order     *OrderRequest `json:"order"`
	Status    OrderStatus   `json:"status"`
	TimeStamp uint
}

type OrderItem struct {
	ItemCode string `json:"item_code"`
	Number   int8   `json:"number"`
}

type OrderData struct {
}
