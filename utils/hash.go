package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/pbkdf2"
	"io"
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

func GeneratePasswordInfo() (newPassword string, newAccountSalt string, newMacSalt string, newSavePassword string, err error) {
	newPass := make([]byte, 4)
	if _, err = rand.Read(newPass); err != nil {
		return
	}
	newPassword = hex.EncodeToString(newPass)
	newAccountSalt = GenerateSalt()
	newMacSalt = GenerateSalt()
	newMasterKey := pbkdf2.Key([]byte(newPassword), []byte(newAccountSalt), 1000, 64, sha512.New)
	hkdfReader := hkdf.New(sha512.New, newMasterKey, []byte{}, []byte("HOME-CLOUD-AUTH-KEY-FOR-LOGIN"))
	newAuth := make([]byte, 32)
	if _, err = io.ReadFull(hkdfReader, newAuth); err != nil {
		return
	}
	newAuthKey := hex.EncodeToString(newAuth)
	newSavePassword = GetHashWithSalt(newAuthKey, newMacSalt)
	return
}
