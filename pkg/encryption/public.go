package encryption

import (
	"crypto/rsa"

	"golang.org/x/crypto/ssh"
)

// GeneratePublicRSAKey generates a public RSA key
func GeneratePublicRSAKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return pubKeyBytes, nil
}
