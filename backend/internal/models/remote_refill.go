package models

import (
	"bar/autogen"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	RemoteRefill struct {
		PrivateID      primitive.ObjectID `bson:"_id,omitempty" json:"-"`
		autogen.RemoteRefill `bson:",inline"`
	}
)

// ToJSON converts the model to JSON
func (o *RemoteRefill) ToJSON() []byte {
	data, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return data
}
