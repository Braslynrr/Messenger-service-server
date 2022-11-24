package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

// EncryptText encrypts a Text using a AES-256 key
func EncryptText(keystr string, stext string) (_ string, err error) {
	key := []byte(keystr)
	text := []byte(stext)

	c, err := aes.NewCipher(key)
	gcm, err := cipher.NewGCM(c)

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	result := gcm.Seal(nonce, nonce, text, nil)

	return string(result), nil
}

// DecryptText decrypts a Text using a AES-256 key
func Decrypt(keystr string, text string) (_ string, err error) {
	ciphertext := []byte(text)
	key := []byte(keystr)
	c, err := aes.NewCipher(key)
	gcm, err := cipher.NewGCM(c)

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext size is less than nonceSize")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)

	return string(plaintext), err
}

//GenerateRandomAESKey generates a Random AES-256 key
func GenerateRandomAESKey() (string, error) {

	key := make([]byte, 16) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(key); err != nil {
		return "", err
	}

	return hex.EncodeToString(key), nil //convert to string for saving

}

func GenerateToken() (string, error) {
	token := make([]byte, 4) //generate a token
	if _, err := rand.Read(token); err != nil {
		return "", err
	}

	return hex.EncodeToString(token), nil //convert to string for saving
}
