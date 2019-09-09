package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"log"
)

// GeneratePrivateRSAKeyPair generates a Private RSA Key and outputs pem file format bytes
func GeneratePrivateRSAKeyPair() (*rsa.PrivateKey, []byte) {
	bitSize := 4096
	privateKey, err := generatePrivateKey(bitSize)
	if err != nil {
		log.Fatal(err.Error())
	}

	privateKeyPEMBytes := encodePrivateKeyToPEM(privateKey)

	return privateKey, privateKeyPEMBytes
}

// PrivateRSADecryptAESKey decrypts the AES key
func PrivateRSADecryptAESKey(privateKey []byte, encryptedText, label []byte) (decryptedText []byte) {
	parsedPrivateKey, perr := x509.ParsePKCS1PrivateKey(privateKey)
	if perr != nil {
		panic(perr)
	}

	hash := sha512.New()
	decryptedText, derr := rsa.DecryptOAEP(hash, rand.Reader, parsedPrivateKey, encryptedText, label)
	if derr != nil {
		panic(derr)
	}

	return decryptedText
}

func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	privateKey, gerr := rsa.GenerateKey(rand.Reader, bitSize)
	if gerr != nil {
		return nil, gerr
	}

	verr := privateKey.Validate()
	if verr != nil {
		return nil, verr
	}

	return privateKey, nil
}

func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}
