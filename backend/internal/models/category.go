package models

import (
	"bar/autogen"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Category struct {
		PrivateID        primitive.ObjectID `bson:"_id,omitempty" json:"-"`
		autogen.Category `bson:",inline"`
	}
)

// ToJSON converts the model to JSON
func (o *Category) ToJSON() []byte {
	data, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return data
}
