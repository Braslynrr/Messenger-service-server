package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

// EncryptInterface encrypts an interface
func EncryptInterface(something any, key string) (encrypted string, err error) {
	var bytes []byte
	var keyInBytes []byte

	bytes, err = json.MarshalIndent(something, "", "")
	if err == nil {
		keyInBytes, err = hex.DecodeString(key)
		if err == nil {
			bytes, err = encrypt_aes_cbc(bytes, keyInBytes)
		}
	}
	base64Text := base64.StdEncoding.EncodeToString(bytes)
	return base64Text, err
}

// EncryptMiddleWare encrypts the body before send it
func EncryptMiddleWare(EncryptedEnabled bool) gin.HandlerFunc {
	return func(ctx *gin.Context) { // pending to fix
		if EncryptedEnabled && ctx.Writer.Status() == http.StatusOK {

			ctx.Done()

			requestBody, err := ioutil.ReadAll(ctx.Request.Response.Body)

			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			session := sessions.Default(ctx)

			key := session.Get("key").(string)

			encryptedBody, err := encrypt_aes_cbc(requestBody, []byte(key))

			newBody := base64.StdEncoding.EncodeToString(encryptedBody)

			ctx.Writer.WriteHeader(http.StatusOK)
			ctx.Writer.Write([]byte{})

			ctx.String(http.StatusOK, newBody)
		}
		ctx.Next()
	}
}

// EncryptMiddleWare encrypts the body before send it
func DecryptMiddleWare(EncryptedEnabled bool) gin.HandlerFunc {
	return func(ctx *gin.Context) { // pending to fix
		if EncryptedEnabled {
			encryptedBody, err := ioutil.ReadAll(ctx.Request.Body)

			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			newBody := base64.StdEncoding.EncodeToString(encryptedBody)

			session := sessions.Default(ctx)

			key := session.Get("key").(string)

			decryptedBody, err := decrypt_aes_cbc([]byte(newBody), []byte(key))

			ctx.Request.Body = io.NopCloser(strings.NewReader(string(decryptedBody)))

		}
		ctx.Next()
	}
}

// hasNumber checks if rune is a digit
func HasNumber(r rune) bool {
	return unicode.IsDigit(r)
}

// filterString filters a string using a function of rune
func FilterString(s string, f func(rune) bool) string {
	var filtered string
	for _, r := range s {
		if f(r) {
			filtered += string(r)
		}
	}
	return filtered
}
