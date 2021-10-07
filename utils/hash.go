package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
)

func GetHash(str string) string {
	hash := sha512.Sum512([]byte(str))
	return hex.EncodeToString(hash[:])
}

func GetHashWithSalt(str string, salt string) string {
	hash := sha512.Sum512(append([]byte(str), []byte(salt)...))
	return hex.EncodeToString(hash[:])
}

func GenerateSalt(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890~!@#$%^&*()_+=-/"
	charLength := len(chars)
	length = length / 4
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		//default salt
		GetLogger().Error("Generate salt error, fallback to use default 128-bit salt")
		return "bjkqjbQWDQ123VWQacaPqlpMokthwCAS"
	}
	for i := 0; i < length; i++ {
		salt[i] = chars[int(salt[i])%charLength]
	}
	return string(salt)
}
