package ssh

import (
	"crypto/rand"
	"crypto/rsa"

	ext_ssh "golang.org/x/crypto/ssh"
)

// KeyPair A SSH Key
type KeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  []byte
}

// NewKeyPair Returns  new SSH Keypain
func NewKeyPair(keySize int) (*KeyPair, error) {
	// generate private RSA key - EC2 instance connect only supports RSA
	rsaKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, err
	}

	// generate public key
	publicKey, err := ext_ssh.NewPublicKey(&rsaKey.PublicKey)
	if err != nil {
		return nil, err
	}

	publicKeySerialized := ext_ssh.MarshalAuthorizedKey(publicKey)

	return &KeyPair{
		PrivateKey: rsaKey,
		PublicKey:  publicKeySerialized,
	}, nil
}
