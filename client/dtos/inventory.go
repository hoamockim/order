package dtos

type Product struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type InventoryList struct {
}

func GetInventoryList() (inventList *InventoryList, err error) {
	return
}
