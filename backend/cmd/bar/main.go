package main

import (
	"bar/api"
	"bar/autogen"
	"bar/internal/config"
	"bar/internal/db"
	"bar/internal/models"
	"context"
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
			CardId:    autogen.OptionalString("1"),
			Id:        uuid.New(),
			Role:      autogen.AccountSuperAdmin,
			State:     autogen.AccountOK,
			PriceRole: autogen.AccountPriceCeten,
		},
	}
	acc.SetPin("1234")

	// Create default user (ignore errors)
	err = db.CreateAccount(context.Background(), acc)
	if err == nil {
		logrus.Infof("Created default user : %#+v", acc)
	}

	s := api.NewServer(db)

	if err := s.Serve(&c); err != nil {
		logrus.Panic(err)
	}
}
