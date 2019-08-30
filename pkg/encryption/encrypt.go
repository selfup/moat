package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
)

// Encrypt will encrypt the text file
func Encrypt(fileContents []byte, aesKey string) []byte {
	key := []byte(aesKey)

	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal("c ", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatal("gcm ", err)
	}

	nonce := make([]byte, gcm.NonceSize())

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		log.Fatal(err)
	}

	return gcm.Seal(nonce, nonce, fileContents, nil)
}
