package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"order/db/models"
	"order/pkg/cache"
	"order/pkg/configs"
	"time"
)

type dbFacade struct {
	app      string
	ctx      context.Context
	database *mongo.Database
	cache    cache.Adapter
}

type Filter struct {
	F []F `bson:"fields"`
}

type F struct {
	Key   string      `json:"key" bson:"key"`
	Value interface{} `json:"value" bson:"value"`
}

var fcd *dbFacade

func InitDB() {
	var (
		client *mongo.Client
		err    error
	)
	if client, err = mongo.NewClient(options.Client().ApplyURI(configs.MongoURI())); err != nil {
		panic(err)
	}

	fcd = new(dbFacade)
	fcd.ctx = context.Background()
	fcd.database = client.Database(configs.MongoDatabase())
	if err = client.Connect(fcd.ctx); err != nil {
		panic(err)
	}
}

func cll(name string) *mongo.Collection {
	return fcd.database.Collection(name)
}

func GetAll(document interface{ models.ModelMetadata }) ([]interface{}, error) {
	fcd.ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	filter := bson.D{{}}
	return filterDocument(filter, document)
}

func GetById(id string, data interface{ models.ModelMetadata }) error {
	fcd.ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}

	if rs := cll(data.CollectionName()).FindOne(fcd.ctx, bson.M{"_id": objectId}); rs != nil {
		err = rs.Decode(data)
	}
	return err
}

func Save(document interface{ models.ModelMetadata }) (err error) {
	fcd.ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	_, err = fcd.database.Collection(document.CollectionName()).InsertOne(fcd.ctx, document)
	return
}

func InsertBatch(collectionName string, documents []interface{}) error {
	fcd.ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
	_, err := cll(collectionName).InsertMany(fcd.ctx, documents)
	return err
}

func Update(collectionName string, documents []interface{}, filters Filter) error {
	fcd.ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	m, _ := makeFilter(filters)
	_, err := cll(collectionName).UpdateOne(fcd.ctx, m, documents)
	return err
}

func CreateCollection(collectionName string) error {
	command := bson.D{{"create", collectionName}}
	var result bson.M
	return fcd.database.RunCommand(context.Background(), command).Decode(&result)
}

func filterDocument(filter interface{}, document interface{ models.ModelMetadata }) ([]interface{}, error) {
	collectionName := document.CollectionName()
	cur, err := cll(collectionName).Find(fcd.ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(fcd.ctx)
	var documents []interface{}
	for cur.Next(fcd.ctx) {
		entity := models.NewModel(collectionName)
		err = cur.Decode(entity)
		if err != nil {
			return nil, err
		}
		documents = append(documents, entity)
	}
	if err = cur.Err(); err != nil {
		return nil, err
	}

	if len(documents) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return documents, nil
}

func makeFilter(filters Filter) (*bson.M, error) {
	m, err := bson.Marshal(&filters)
	if err != nil {
		return nil, err
	}
	var doc bson.M
	_ = bson.Unmarshal(m, &doc)
	return &doc, nil
}
