package dto

type OrderRequest struct {
	Items       []OrderItem `json:"items"`
	CustomerId  string      `json:"customer_id"`
	AddressCode string      `json:"address_code"`
	CreatedAt   int         `json:"created_at"`
}

type OrderItem struct {
	ItemCode string `json:"item_code"`
	Number   int8   `json:"number"`
}

type OrderData struct {
}
