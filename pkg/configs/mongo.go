package configs

type Mongo struct {
	Url string `json:"url"`
	Database string `json:"database" env:"MONGO_DB"`
}

func (db *Mongo) isValid() bool {
	return true
}

func MongoURI() string {
	return app.Mongo.Url
}

func MongoDatabase() string {
	return app.Mongo.Database
}
