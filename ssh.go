package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// SSHKeyPair A SSH Key
type SSHKeyPair struct {
	private *rsa.PrivateKey
	public  []byte
}

// NewSSHKeyPair Returns  new SSH Keypain
func NewSSHKeyPair(keySize int) (*SSHKeyPair, error) {
	// generate private RSA key - EC2 instance connect only supports RSA
	rsaKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, err
	}

	// generate public key
	publicKey, err := ssh.NewPublicKey(&rsaKey.PublicKey)
	if err != nil {
		return nil, err
	}

	publicKeySerialized := ssh.MarshalAuthorizedKey(publicKey)

	return &SSHKeyPair{
		private: rsaKey,
		public:  publicKeySerialized,
	}, nil
}

// SSHAgent A SSH agent
type SSHAgent struct {
	conn net.Conn
}

// NewSSHAgent Make a new coonection to the SSH Agent
func NewSSHAgent() (*SSHAgent, error) {
	socket := os.Getenv("SSH_AUTH_SOCK")
	connection, err := net.Dial("unix", socket)
	if err != nil {
		log.Fatalf("failed to open SSH_AUTH_SOCK: %v", err)
	}
	return &SSHAgent{
		conn: connection,
	}, nil
}

func (s SSHAgent) addKey(key *rsa.PrivateKey, comment string) error {
	log.Debug("adding key to agent")
	agentClient := agent.NewClient(s.conn)
	tmpKey := agent.AddedKey{
		PrivateKey:       key,
		Comment:          fmt.Sprintf("essh:%s", comment),
		LifetimeSecs:     10,
		ConfirmBeforeUse: false,
	}

	err := agentClient.Add(tmpKey)
	if err != nil {
		log.Fatal("could not add key to ssh agent")
		return err
	}

	return nil
}

func sshConnect(args []string) error {
	cmd := exec.Command("ssh", args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
