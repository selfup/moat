package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"

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

// PublicRSAEncryptAESKey encrypts the AES key
func PublicRSAEncryptAESKey(publicKey *rsa.PublicKey, sourceText, label []byte) (encryptedText []byte) {
	hash := sha512.New()

	encryptedText, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, sourceText, label)
	if err != nil {
		panic(err)
	}

	return encryptedText
}
