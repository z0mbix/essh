package ssh

import (
	"crypto/rand"
	"crypto/rsa"

	extssh "golang.org/x/crypto/ssh"
)

// KeyPair A SSH Key
type KeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  []byte
}

// NewKeyPair Returns a new RSA SSH Keypair
func NewKeyPair(keySize int) (*KeyPair, error) {
	// generate private RSA key - EC2 instance connect only supports RSA
	rsaKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, err
	}

	// generate public key
	publicKey, err := extssh.NewPublicKey(&rsaKey.PublicKey)
	if err != nil {
		return nil, err
	}

	publicKeySerialized := extssh.MarshalAuthorizedKey(publicKey)

	return &KeyPair{
		PrivateKey: rsaKey,
		PublicKey:  publicKeySerialized,
	}, nil
}
