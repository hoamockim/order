package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseModel struct {
	CreatedAt uint
	UpdatedAt uint
}

type ModelMetadata interface {
	CollectionName() string
}

type ModelUpdate interface {
	MakeCondition(filterFields []string) map[string]interface{}
	MakeUpdated(fields []string) []primitive.E
}

func NewModel(collectionName string) interface{ ModelMetadata } {
	switch collectionName {
	case "order":
		return &OrderInfo{}
	default:
		return nil
	}
}
