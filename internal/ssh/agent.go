package ssh

import (
	"crypto/rsa"
	"fmt"
	"net"
	"os"

	"github.com/apex/log"
	"golang.org/x/crypto/ssh/agent"
)

// Agent is an SSH agent
type Agent struct {
	conn net.Conn
}

// NewAgent Make a new connection to the SSH Agent
func NewAgent() (*Agent, error) {
	socket := os.Getenv("SSH_AUTH_SOCK")
	connection, err := net.Dial("unix", socket)
	if err != nil {
		return nil, fmt.Errorf("failed to open SSH_AUTH_SOCK: %s", err)
	}
	return &Agent{
		conn: connection,
	}, nil
}

func (s Agent) addKey(key *rsa.PrivateKey, comment string, timeout uint32) error {
	log.Debugf("adding key to agent for %d seconds", timeout)
	agentClient := agent.NewClient(s.conn)
	tmpKey := agent.AddedKey{
		PrivateKey:       key,
		Comment:          fmt.Sprintf("essh:%s", comment),
		LifetimeSecs:     timeout,
		ConfirmBeforeUse: false,
	}

	err := agentClient.Add(tmpKey)
	if err != nil {
		return fmt.Errorf("could not add key to ssh agent: %s", err)
	}

	return nil
}
