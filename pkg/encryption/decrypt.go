package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"log"
)

// Decrypt decrypts the text file
func Decrypt(contents []byte, aesKey string) []byte {
	key := []byte(aesKey)

	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal("decrypt newcypher ", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatal("decrypt newgcm ", err)
	}

	nonceSize := gcm.NonceSize()
	nonce, stuff := contents[:nonceSize], contents[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, stuff, nil)
	if err != nil {
		log.Fatal("decrypt open ", err)
	}

	return plaintext
}
