package ssh

import (
	"crypto/rsa"
	"fmt"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
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
		log.Fatalf("failed to open SSH_AUTH_SOCK: %v", err)
	}
	return &Agent{
		conn: connection,
	}, nil
}

func (s Agent) addKey(key *rsa.PrivateKey, comment string) error {
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
