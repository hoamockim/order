package models

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type BaseModel struct {
	Id        int
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type ModelCredential interface {
	Validate() bool
}

type ModelMetadata interface {
	GetDocumentName() string
}

type BsonMod interface {
	GenerateBson() bson.D
}

type ModelCache interface {
	IsCached() bool
}

const (
	DataInvalid string = "Data is invalid"
)
