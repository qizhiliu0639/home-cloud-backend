package utils

import (
	"crypto/rand"
	"crypto/sha256"
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

func GenerateSalt() string {
	// 256-bit salt
	length := 32
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		//default salt
		GetLogger().Error("Generate salt error, fallback to use default 256-bit salt")
		return "166845ab354965a468de3bce654f1199294812bbddcceed46986bcbc7823ccad"
	}
	return hex.EncodeToString(salt[:])
}

func GenerateFakeSalt(username string) string {
	hash := sha256.Sum256(append([]byte(username), []byte("FQMMWDqwsdq@!234DFQAWASCASEDQOAOS@#$#)T!$(@#")...))
	return hex.EncodeToString(hash[:])
}
