package main

import (
	"bar/internal/config"
	"bar/internal/db"
	"time"

	"bar/internal/db/mongo"

	"github.com/sirupsen/logrus"
)

func main() {
	c := config.GetConfig()

	opts := db.NewDatabaseOptions(c.MongoConfig.ConnectionURI, c.MongoConfig.Database, time.Millisecond*time.Duration(c.MongoConfig.Timeout))
	db := mongo.NewMongoBackend(opts)

	err := db.Connect()
	if err != nil {
		logrus.Panic(err)
	}

}
