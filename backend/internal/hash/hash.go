package hash

import (
	"github.com/matthewhartstonge/argon2"
)

var argon = argon2.DefaultConfig()

func Hash(pwd string) (string, error) {
	encoded, err := argon.HashEncoded([]byte(pwd))
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

func MustHash(pwd string) string {
	hash, err := Hash(pwd)
	if err != nil {
		panic(err)
	}

	return hash
}

func Verify(hash, pwd string) (bool, error) {
	ok, err := argon2.VerifyEncoded([]byte(pwd), []byte(hash))
	if err != nil {
		return false, err
	}

	return ok, nil
}

func MustVerify(hash, pwd string) bool {
	ok, err := Verify(hash, pwd)
	if err != nil {
		panic(err)
	}

	return ok
}
