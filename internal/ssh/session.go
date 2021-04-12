package ssh

import (
	"os"
	"os/exec"
	"strings"

	"github.com/apex/log"

	"github.com/z0mbix/essh/internal/aws"
	"github.com/z0mbix/essh/internal/config"
)

// Session wraps all ssh actions
type Session struct {
	KeyPair *KeyPair
	Agent   *Agent
}

// NewSession creates a new ssh session with a new
// ssh keypair created and put into the ssh agent
func NewSession(sessionID string, timeout uint32) (*Session, error) {
	keypair, err := NewKeyPair(2048)
	if err != nil {
		return nil, err
	}

	log.Debug("creating new ssh agent connection")
	agent, err := NewAgent()
	if err != nil {
		return nil, err
	}

	log.Debug("adding private key to ssh agent")
	if err = agent.addKey(keypair.PrivateKey, sessionID, timeout); err != nil {
		return nil, err
	}

	return &Session{KeyPair: keypair, Agent: agent}, nil
}

// Connect executes ssh command to the given instance with args passed to the command line
func (s *Session) Connect(instance *aws.Instance, args []string) error {
	if err := instance.SendPublicKey(config.UserName, string(s.KeyPair.PublicKey)); err != nil {
		return err
	}

	log.Infof("running command: ssh %s\n", strings.Join(args[:], " "))

	cmd := exec.Command("ssh", args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
