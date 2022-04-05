package repositories

import (
	"order/db"
	"order/db/models"
)

type defaultRepository struct {
}

//FilterField is used for where clause
type FilterField struct {
	Column    string
	Value     interface{}
	And       bool
	Or        bool
	condition string // =, !=, >=, <=, >, <
}

var defaultRepo defaultRepository

func New() *defaultRepository {
	return &defaultRepo
}

//save
func (repo *defaultRepository) save(data interface {
	models.ModelCredential
	models.ModelMetadata
	models.BsonMod
}) error {
	return db.Save(data.GetDocumentName(), data.GenerateBson())
}

//getById
func (repo *defaultRepository) getById(data interface {
	models.ModelMetadata
	models.ModelCache
}, id string) error {
	return db.GetById(id, data.GetDocumentName(), data)
}

//update
func (repo *defaultRepository) update(data interface {
	models.ModelCredential
	models.ModelMetadata
}, filter ...FilterField) error {
	var conditions map[string]interface{}
	if len(filter) > 0 {
		for _, val := range filter {
			key := val.Column
			conditions[key] = val.Value
		}
	}
	return db.Update(data.GetDocumentName(), data, conditions)
}

//filter
func (repo *defaultRepository) filter(data interface {
	models.ModelMetadata
	models.ModelCache
}, entities interface {}, filter ...FilterField) error {
	var conditions map[string]interface{}
	if len(filter) > 0 {
		for _, val := range filter {
			key := val.Column
			if val.Or {
			}
			conditions[key] = val.Value
		}
	}
	return db.Filter(data.GetDocumentName(), entities, conditions)
}

