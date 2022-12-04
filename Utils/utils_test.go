package utils_test

import (
	"MessengerService/utils"
	"encoding/json"

	"testing"
)

const MESSAGE = "◉"

// TestEncryptAndDecrypy calls EncryptText and DecryptText checking
func TestEncryptAndDecryt(t *testing.T) {
	key, err := utils.GenerateRandomAESKey(16)

	if err != nil {
		t.Fatalf("Generating a key should not return an error. Error: %v", err)
	}

	encryptedMessage, err := utils.EncryptText(MESSAGE, key)
	if err != nil {
		t.Fatalf("Encrypting the message should not return an error. Error: %v", err)
	}

	decryptedMessage, err := utils.DecryptText(encryptedMessage, key)

	if err != nil {
		t.Fatalf("Decrypting the message should not return an error. Error: %v", err)
	}

	if MESSAGE != string(decryptedMessage) {
		t.Fatalf("Decrypted Message should be equal to original message. Original message: %v != Decrypted Message: %v", MESSAGE, string(encryptedMessage))
	}
}

// TestGenerateDifferentKeys calls GenerateRandomAESKey checking all outputs are different
func TestGenerateDifferentKeys(t *testing.T) {
	generatedKeys := make(map[string]bool, 0)

	for i := 0; i < 10; i++ {
		key, err := utils.GenerateRandomAESKey(16)
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

// TestEncryptAndDecrytInterface calls EncryptInterface to encrypt and interface
func TestEncryptAndDecrytInterface(t *testing.T) {
	type example struct {
		Example   string
		Something int
	}
	key, _ := utils.GenerateRandomAESKey(16)
	exampleVariable := example{Example: "example", Something: 1}

	encryptedText, err := utils.EncryptInterface(exampleVariable, key)

	if err != nil {
		t.Fatalf("EncryptInterface should not return an error: %v", err.Error())
	}

	plain, err := utils.DecryptText(encryptedText, key)
	if err != nil {
		t.Fatalf("DecryptInterface should not return an error: %v", err.Error())
	}

	err = json.Unmarshal([]byte(plain), &exampleVariable)
	if err != nil {
		t.Fatalf("Unmarshal should be ok. erro: %v", err)
	}
}
