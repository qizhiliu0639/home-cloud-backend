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
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		//Todo Logger Set
		return ""
	}
	for i := 0; i < length; i++ {
		salt[i] = chars[int(salt[i])%charLength]
	}
	return string(salt)
}
