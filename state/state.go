package state

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"
)

// EncryptState - Creates an encrypted state for later verification
func EncryptState(clientID string, key string) string {
	hashedKey := createHash(addSalt(key))
	encryptedState, err := encrypt([]byte(clientID), []byte(hashedKey))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(encryptedState)
}

// DecryptState - Decrypts the given state and returns the clientID or an error
func DecryptState(state string, key string) (string, error) {
	encryptedState, err := base64.StdEncoding.DecodeString(state)
	if err != nil {
		return "", err
	}

	hashedKey := createHash(addSalt(key))
	clientID, err := decrypt([]byte(encryptedState), []byte(hashedKey))
	if err != nil {
		return "", err
	}
	return string(clientID), nil
}

func addSalt(key string) string {
	salt := time.Now().Format("2006-01-02")
	return fmt.Sprintf("%s-%s", key, salt)
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
