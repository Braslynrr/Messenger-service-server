package utils_test

import (
	"MessengerService/utils"
	"testing"
)

const MESSAGE = "I Will be Encrypted"

// TestEncryptAndDecrypy calls EncryptText and Decrypt checking
func TestEncryptAndDecrypy(t *testing.T) {

	keyStr, err := utils.GenerateRandomAESKey()
	if err != nil {
		t.Fatalf("Generating a key should not return an error. Error: %v", err)
	}
	encryptedMessage, err := utils.EncryptText(keyStr, MESSAGE)
	if err != nil {
		t.Fatalf("Encrypting the message should not return an error. Error: %v", err)
	}

	decryptedMessage, err := utils.Decrypt(keyStr, encryptedMessage)
	if err != nil {
		t.Fatalf("Decrypting the message should not return an error. Error: %v", err)
	}

	if MESSAGE != decryptedMessage {
		t.Fatalf("Decrypted Message should be equal to original message. Original message: %v != Decrypted Message: %v", MESSAGE, encryptedMessage)
	}
}

// TestGenerateDifferentKeys calls GenerateRandomAESKey checking all outputs are different
func TestGenerateDifferentKeys(t *testing.T) {
	generatedKeys := make(map[string]bool, 0)

	for i := 0; i < 10; i++ {
		key, err := utils.GenerateRandomAESKey()
		if err != nil {
			t.Fatalf("Generating a key1 should not return an error. Error: %v", err)
		}
		if !generatedKeys[key] {
			generatedKeys[key] = true
		} else {
			t.Fatalf("Generated keys should not have this key: %v", key)
		}
	}

}
