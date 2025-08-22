package models

import (
	"bar/autogen"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Starring struct {
		PrivateID      primitive.ObjectID `bson:"_id,omitempty" json:"-"`
		autogen.Starring `bson:",inline"`
	}
)

// ToJSON converts the model to JSON
func (o *Starring) ToJSON() []byte {
	data, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return data
}
