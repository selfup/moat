package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"hash"
	"log"

	"golang.org/x/crypto/ssh"
)

// GeneratePublicRSAKey generates a public RSA key
func GeneratePublicRSAKey(publicKey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return pubKeyBytes, nil
}

// EncryptAESKey encrypts the AES key
func EncryptAESKey(publicKey *rsa.PublicKey, sourceText, label []byte) (encryptedText []byte) {
	var err error
	var hash hash.Hash
	hash = sha512.New()
	if encryptedText, err = rsa.EncryptOAEP(hash, rand.Reader, publicKey, sourceText, label); err != nil {
		log.Fatal(err)
	}

	return encryptedText
}
