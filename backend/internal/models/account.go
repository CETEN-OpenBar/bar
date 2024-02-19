package models

import (
	"bar/autogen"
	"bar/internal/hash"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Account struct {
		PrivateID       primitive.ObjectID `bson:"_id,omitempty" json:"-"`
		autogen.Account `bson:",inline"`

		CreatedAt int64 `bson:"created_at" json:"created_at"`
	}
)

func (o *Account) SetPin(pwd string) {
	o.Account.CardPin = hash.MustHash(pwd)
}

func (o *Account) SetPassword(pwd string) {
	h := hash.MustHash(pwd)
	o.Account.Password = &h
}

func (o *Account) VerifyPin(pwd string) bool {
	ok, _ := hash.Verify(o.Account.CardPin, pwd)
	return ok
}

func (o *Account) VerifyPassword(pwd string) bool {
	ok, _ := hash.Verify(*o.Account.Password, pwd)
	return ok
}

func (o *Account) IsAdmin() bool {
	return o.Role == autogen.AccountAdmin || o.Role == autogen.AccountSuperAdmin || o.Role == autogen.AccountMember
}

func (o *Account) IsBlocked() bool {
	var blocked bool
	for _, res := range o.Restrictions {
		if res == autogen.AccountBlocked {
			blocked = true
		}
	}
	return blocked
}

// ToJSON converts the model to JSON
func (o *Account) ToJSON() []byte {
	data, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return data
}
