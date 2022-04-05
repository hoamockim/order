package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"order/pkg/cache"
	"order/pkg/configs"
	"time"
)

const (
	DBDriver = "mongoDB"
)

type dbFacade struct {
	app      string
	ctx      context.Context
	database *mongo.Database
	cache    cache.Adapter
}

var fcd *dbFacade

func init() {
	var (
		client *mongo.Client
		err    error
	)

	if client, err = mongo.NewClient(options.Client().ApplyURI(configs.MongoURI())); err != nil {
		panic(err)
	}
	if err = client.Connect(fcd.ctx); err != nil {
		panic(err)
	}

	fcd = new(dbFacade)
	fcd.ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	fcd.database = client.Database(configs.MongoDatabase())
}

func GetById(id, document string, data interface{}) (err error) {
	if rs := fcd.database.Collection(document).FindOne(fcd.ctx, bson.D{{"_id", id}}); rs != nil {
		err = rs.Decode(data)
	}
	return
}

func Update(tableName string, data interface{}, conditions map[string]interface{}) error {
	return nil
}

func Filter(tableName string, entities interface{}, conditions map[string]interface{}) error {
	return nil
}

func Save(document string, doc bson.D) (err error) {
	_, err = fcd.database.Collection(document).InsertOne(fcd.ctx, doc)
	return
}

func SaveMany(document string, docs []interface{}) (err error) {
	_, err = fcd.database.Collection(document).InsertMany(fcd.ctx, docs)
	return
}
