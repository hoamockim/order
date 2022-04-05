package dtos

type TransactionRequest struct {
	Filter   *TransactionFilter `json:"filter"`
	Pageable struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	} `json:"pageable"`
	Sort struct {
		TransactionCreateTime string `json:"transaction_create_time"`
	}
}

type TransactionFilter struct {
	ActivityType struct {
		In []string `json:"in"`
	} `json:"activity_type"`
	GroupId struct {
		In []string `json:"in"`
	} `json:"group_id"`
	TransactionStatus struct {
		In []string `json:"in"`
	} `json:"transaction_status"`
}

type TransactionResponse struct {
	Meta *Meta              `json:"meta,omitempty"`
	Data []*TransactionData `json:"data"`
}

type Meta struct {
	Code   int `json:"code"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type TransactionData struct {
	TransactionId            string     `json:"transaction_id"`
	TransactionRef           string     `json:"transaction_ref"`
	TransactionAmount        float64    `json:"transaction_amount"`
	TransactionAmountCharged float64    `json:"transaction_amount_charged"`
	TransactionStatus        string     `json:"transaction_status"`
	ActivityType             string     `json:"activity_type"`
	SourceType               string     `json:"source_type"`
	DestinationType          string     `json:"destination_type"`
	GroupId                  string     `json:"group_id"`
	CustomerName             string     `json:"customer_name,omitempty"`
	TransactionUpdateTime    int64      `json:"transaction_update_time"`
	TransactionCreateTime    int64      `json:"transaction_create_time"`
	Description              string     `json:"description"`
	ExtraData                *ExtraData `json:"extra_data,omitempty"`
}

type ExtraData struct {
	Description string `json:"description,omitempty"`
}

