package models

import (
	"bar/autogen"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	CashMovement struct {
		PrivateID            primitive.ObjectID `bson:"_id,omitempty" json:"-"`
		autogen.CashMovement `bson:",inline"`
	}
)

// ToJSON converts the model to JSON
func (o *CashMovement) ToJSON() []byte {
	data, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return data
}
