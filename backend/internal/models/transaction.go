package models

import (
	"bar/autogen"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Transaction struct {
		PrivateID primitive.ObjectID `bson:"_id,omitempty" json:"-"`
		autogen.Transaction
	}
)

// ToJSON converts the model to JSON
func (o *Transaction) ToJSON() []byte {
	data, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return data
}
