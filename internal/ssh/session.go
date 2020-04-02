package ssh

import (
	"os"
	"os/exec"
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

// Connect executes ssh command with a given args
func (s *Session) Connect(args []string) error {
	cmd := exec.Command("ssh", args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
