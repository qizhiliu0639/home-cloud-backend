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

var fakeSaltSalt = GenerateSaltOrKey()

// GetHashWithSalt Assume str and salt are encoded to hex, if not will return empty string (will not match)
func GetHashWithSalt(str string, salt string) string {
	strBytes, err := hex.DecodeString(str)
	if err != nil {
		return ""
	}
	var saltBytes []byte
	saltBytes, err = hex.DecodeString(salt)
	if err != nil {
		return ""
	}
	hash := sha512.Sum512(append(strBytes, saltBytes...))
	return hex.EncodeToString(hash[:])
}

// GenerateSaltOrKey generate a random 256-bit for encryption or hash, in hex format
func GenerateSaltOrKey() string {
	// 256-bit salt or key for encryption
	length := 32
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		//default salt
		GetLogger().Error("Generate salt error, fallback to use default 256-bit salt")
		return "166845ab354965a468de3bce654f1199294812bbddcceed46986bcbc7823ccad"
	}
	return hex.EncodeToString(salt[:])
}

// GenerateFakeSalt generate a fake salt using SHA-256 with not exist username
// and a random salt generate at program starts
func GenerateFakeSalt(username string) string {
	fss, err := hex.DecodeString(fakeSaltSalt)
	if err != nil {
		// if error occur, use the fall back salt
		fss = []byte("FQMMWDqwsdq@!234DFQAWASCASEDQOAOS@#$#)T!$(@#")
	}
	hash := sha256.Sum256(append([]byte(username), fss...))
	return hex.EncodeToString(hash[:])
}

// GeneratePasswordInfo Generate authentication information for a new account
// newAccountSalt is encoded in hex and decoded directly
// because of compatibility with frontend salt generation
func GeneratePasswordInfo() (
	newPassword string,
	newAccountSalt string,
	newMacSalt string,
	newSavePassword string,
	newEncryptionKey string,
	err error,
) {
	newPass := make([]byte, 4)
	if _, err = rand.Read(newPass); err != nil {
		return
	}
	newPassword = hex.EncodeToString(newPass)
	newAccountSalt, newMacSalt, newSavePassword, newEncryptionKey, err = GeneratePasswordInfoFromPassword(newPassword)
	return
}

// GeneratePasswordInfoFromPassword Generate authentication information for account with provided password
// newAccountSalt is encoded in hex and decoded directly
// because of compatibility with frontend salt generation
func GeneratePasswordInfoFromPassword(newPassword string) (
	newAccountSalt string,
	newMacSalt string,
	newSavePassword string,
	newEncryptionKey string,
	err error,
) {
	newAccountSalt = GenerateSaltOrKey()
	newMacSalt = GenerateSaltOrKey()
	newMasterKey := pbkdf2.Key([]byte(newPassword), []byte(newAccountSalt), 1000, 64, sha512.New)
	hkdfReader := hkdf.New(sha512.New, newMasterKey, []byte{}, []byte("HOME-CLOUD-AUTH-KEY-FOR-LOGIN"))
	newAuth := make([]byte, 32)
	if _, err = io.ReadFull(hkdfReader, newAuth); err != nil {
		return "", "", "", "", err
	}
	newAuthKey := hex.EncodeToString(newAuth)
	newSavePassword = GetHashWithSalt(newAuthKey, newMacSalt)
	hkdfReader = hkdf.New(sha512.New, newMasterKey, []byte{}, []byte("HOME-CLOUD-ENCRYPTION-KEY-FOR-FILES"))
	newEncryptKey := make([]byte, 32)
	if _, err = io.ReadFull(hkdfReader, newEncryptKey); err != nil {
		return "", "", "", "", err
	}
	var encryptKey []byte
	encryptKey, err = hex.DecodeString(GenerateSaltOrKey())
	if err != nil {
		return "", "", "", "", err
	}
	newEncryptionKey, err = EncryptEncryptionKey(newEncryptKey, encryptKey)
	if err != nil {
		return "", "", "", "", err
	}
	return
}
