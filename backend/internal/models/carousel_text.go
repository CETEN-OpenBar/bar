package models

import (
	"bar/autogen"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	CarouselText struct {
		PrivateID            primitive.ObjectID `bson:"_id,omitempty" json:"-"`
		autogen.CarouselText `bson:",inline"`

		CreatedAt int64 `bson:"created_at" json:"created_at"`
	}
)

// ToJSON converts the model to JSON
func (o *CarouselText) ToJSON() []byte {
	data, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return data
}
