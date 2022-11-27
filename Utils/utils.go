package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"

	"github.com/zenazn/pkcs7pad"
)

const iv = "fb4c5e213749eddadf1e22d723eaf207"

var IV, _ = hex.DecodeString(iv)

// encrypt_aes_cbc encrypts AES CBC with pkcs7 padding
func encrypt_aes_cbc(plain, key []byte) (encrypted []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return encrypted, err
	}
	blockSize := block.BlockSize()
	origData := pkcs7pad.Pad(plain, blockSize)

	encrypted = make([]byte, len(origData))
	stream := cipher.NewCBCEncrypter(block, IV)
	stream.CryptBlocks(encrypted, origData)
	return
}

// decrypt_aes_cbc decrypt aes cbc with pkcs7 padding
func decrypt_aes_cbc(encrypted, key []byte) (origData []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return encrypted, err
	}
	plain := make([]byte, len(encrypted))
	stream := cipher.NewCBCDecrypter(block, IV)
	stream.CryptBlocks(plain, encrypted)
	origData, err = pkcs7pad.Unpad(plain)
	return
}

// GenerateRandomAESKey generates a Random AES key
func GenerateRandomAESKey(size int) (string, error) {

	key := make([]byte, size) //generate a random key for AES
	if _, err := rand.Read(key); err != nil {
		return "", err
	}

	return hex.EncodeToString(key), nil //convert to string for saving
}

// GenerateToken generates a random token
func GenerateToken() (string, error) {
	token := make([]byte, 4) //generate a token
	if _, err := rand.Read(token); err != nil {
		return "", err
	}

	return hex.EncodeToString(token), nil //convert to string for saving
}

// DecryptText decrypts text using AES CBC PKC7
func DecryptText(encodedtext string, key string) (string, error) {
	keyInBytes, err := hex.DecodeString(key)
	if err != nil {
		return "", err
	}
	text, err := base64.StdEncoding.DecodeString(encodedtext)
	if err != nil {
		return "", err
	}
	plain, err := decrypt_aes_cbc(text, keyInBytes)
	return string(plain), err
}

// EncryptText Encrypt text using AES CBC PKC7
func EncryptText(text string, key string) (string, error) {
	keyInBytes, err := hex.DecodeString(key)
	if err != nil {
		return "", err
	}
	plaintext := []byte(text)
	encryptedText, err := encrypt_aes_cbc(plaintext, keyInBytes)
	base64Text := base64.StdEncoding.EncodeToString(encryptedText)
	return base64Text, err
}
