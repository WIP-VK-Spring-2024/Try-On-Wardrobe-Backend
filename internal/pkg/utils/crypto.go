package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"

	"golang.org/x/crypto/argon2"
)

const saltSize = 5

func Hash(pass, salt []byte) []byte {
	bytes := argon2.IDKey(pass, salt, 1, 64*1024, 4, 32)
	result := make([]byte, base64.StdEncoding.EncodedLen(len(bytes)))
	base64.StdEncoding.Encode(result, bytes)
	return result
}

func NewSalt() ([]byte, error) {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	hexEncoded := make([]byte, hex.EncodedLen(saltSize))
	hex.Encode(hexEncoded, salt)

	return hexEncoded, nil
}
