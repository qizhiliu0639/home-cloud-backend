package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/chacha20poly1305"
)

// This file contains the wrapper functions for encryption and decryption in AEAD mode

// generateNonce generate the 12-byte nonce for AES and ChaCha20-Poly1305
func generateNonce() ([]byte, error) {
	nonce := make([]byte, 12)
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	return nonce, nil
}

// generateNonceForXChaCha generate the 24-byte nonce for XChaCha20-Poly1305
func generateNonceForXChaCha() ([]byte, error) {
	nonce := make([]byte, 24)
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	return nonce, nil
}

// EncryptEncryptionKey used to encrypt the encryption by the encryption key derived from user password
// return hex encode string with nonce and cipher text
func EncryptEncryptionKey(key []byte, encryptionKey []byte) (string, error) {
	var block cipher.Block
	var gcm cipher.AEAD
	var err error
	block, err = aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err = cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	var nonce []byte
	nonce, err = generateNonce()
	if err != nil {
		return "", err
	}
	cipherText := gcm.Seal(nil, nonce, encryptionKey, nil)
	return hex.EncodeToString(append(nonce, cipherText...)), nil
}

// DecryptEncryptionKey decrypt the encrypted key above
func DecryptEncryptionKey(key []byte, encrypted string) ([]byte, error) {
	var block cipher.Block
	var gcm cipher.AEAD
	var err error
	block, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err = cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	var encryptedBytes []byte
	encryptedBytes, err = hex.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}
	var plainText []byte
	plainText, err = gcm.Open(nil, encryptedBytes[:12], encryptedBytes[12:], nil)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

// EncryptFileAES used to encrypt uploaded files in AES
func EncryptFileAES(key []byte, fileContent []byte) ([]byte, error) {
	var block cipher.Block
	var gcm cipher.AEAD
	var err error
	block, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err = cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	var nonce []byte
	nonce, err = generateNonce()
	if err != nil {
		return nil, err
	}
	cipherText := gcm.Seal(nil, nonce, fileContent, nil)
	return append(nonce, cipherText...), nil
}

// DecryptFileAES used to decrypt files in AES
func DecryptFileAES(key []byte, encryptedContent []byte) (plainText []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("error decryption")
		}
	}()
	var block cipher.Block
	var gcm cipher.AEAD
	block, err = aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err = cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plainText, err = gcm.Open(nil, encryptedContent[:12], encryptedContent[12:], nil)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

// EncryptFileChaCha used to encrypt uploaded files in ChaCha20-Poly1305
func EncryptFileChaCha(key []byte, fileContent []byte) ([]byte, error) {
	chaCha20, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, err
	}
	var nonce []byte
	nonce, err = generateNonce()
	if err != nil {
		return nil, err
	}
	cipherText := chaCha20.Seal(nil, nonce, fileContent, nil)
	return append(nonce, cipherText...), nil
}

// DecryptFileChaCha used to decrypted files in ChaCha20-Poly1305
func DecryptFileChaCha(key []byte, encryptedContent []byte) (plainText []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("error decryption")
		}
	}()
	var chaCha20 cipher.AEAD
	chaCha20, err = chacha20poly1305.New(key)
	if err != nil {
		return nil, err
	}
	plainText, err = chaCha20.Open(nil, encryptedContent[:12], encryptedContent[12:], nil)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

// EncryptFileXChaCha used to encrypt uploaded files in XChaCha20-Poly1305
func EncryptFileXChaCha(key []byte, fileContent []byte) ([]byte, error) {
	xChaCha20, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}
	var nonce []byte
	nonce, err = generateNonceForXChaCha()
	if err != nil {
		return nil, err
	}
	cipherText := xChaCha20.Seal(nil, nonce, fileContent, nil)
	return append(nonce, cipherText...), nil
}

// DecryptFileXChaCha used to decrypted files in XChaCha20-Poly1305
func DecryptFileXChaCha(key []byte, encryptedContent []byte) (plainText []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("error decryption")
		}
	}()
	var xChaCha20 cipher.AEAD
	xChaCha20, err = chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}
	plainText, err = xChaCha20.Open(nil, encryptedContent[:24], encryptedContent[24:], nil)
	if err != nil {
		return nil, err
	}
	return plainText, nil
}
