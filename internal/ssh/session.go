package ssh

import (
	log "github.com/sirupsen/logrus"
	"github.com/z0mbix/essh/internal/aws"
	"github.com/z0mbix/essh/internal/config"
	"os"
	"os/exec"
	"strings"
)

// Session wraps all ssh actions
type Session struct {
	KeyPair *KeyPair
	Agent   *Agent
}

// NewSession creates a new ssh session with a new
// ssh keypair created and put into the ssh agent
func NewSession(sessionID string) (*Session, error) {
	keypair, err := NewKeyPair(2048)
	if err != nil {
		return nil, err
	}

	agent, err := NewAgent()
	if err != nil {
		return nil, err
	}

	if err = agent.addKey(keypair.PrivateKey, sessionID); err != nil {
		return nil, err
	}

	return &Session{KeyPair: keypair, Agent: agent}, nil
}

// Connect executes ssh command to the given instance with args passed to the command line
func (s *Session) Connect(instance *aws.Instance, args []string) error {
	if err := instance.SendPublicKey(config.UserName, string(s.KeyPair.PublicKey)); err != nil {
		return err
	}

	log.Printf("running command: ssh %s", strings.Join(args[:], " "))

	cmd := exec.Command("ssh", args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
