package main

import (
	"bar/api"
	"bar/autogen"
	"bar/internal/config"
	"bar/internal/db"
	"bar/internal/models"
	"time"

	"bar/internal/db/mongo"

	"github.com/google/uuid"
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

	acc := &models.Account{
		Account: autogen.Account{
			CardId:  "1",
			CardPin: "03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4",
			Id:      uuid.New(),
			Role:    autogen.AccountSuperAdmin,
			State:   autogen.AccountOk,
		},
	}

	// Create default user (ignore errors)
	err = db.CreateAccount(acc)
	if err == nil {
		logrus.Infof("Created default user : %#+v", acc)
	}

	s := api.NewServer(db)

	if err := s.Serve(&c); err != nil {
		logrus.Panic(err)
	}
}
