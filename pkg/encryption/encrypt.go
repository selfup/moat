package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
)

// Encrypt will encrypt the text file
func Encrypt(fileContents []byte, aesKey string) []byte {
	key := []byte(aesKey)

	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal("encrypt newcypher ", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatal("encrypt newgcm ", err)
	}

	nonce := make([]byte, gcm.NonceSize())

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		log.Fatal("encrypt readfull ", err)
	}

	return gcm.Seal(nonce, nonce, fileContents, nil)
}
