package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetLevel(log.InfoLevel)

	log.SetFormatter(&log.TextFormatter{
		TimestampFormat:        "2006-01-02T15:04:05.000",
		DisableTimestamp:       true,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})

	config, err := getESSHConfig()

	if err != nil {
		fmt.Println("Failed to parse config...")
	}

	if config.Debug {
		log.Debug(spew.Sdump(config))
	}

	var instanceID string

	sshArgs := []string{"-l", config.UserName}
	sshExtraArgs := config.sshExtraArgs

	sess, err := NewAwsSession(config.Region)
	if err != nil {
		log.Fatalf("could not get instance/session: %s", err)
	}

	if config.SearchMode == SearchModeTag {
		log.Debugf("using Name tag %s to find instance id", config.SearchValue)

		//TODO: change this to return more than one result, then show a menu for selection
		instanceID, err = getInstanceIDFromNameTag(sess, config.SearchValue)
		if err != nil {
			log.Fatal(err)
		}
		log.Debugf("found instance id: %s", instanceID)
	} else if config.SearchMode == SearchModeInst {
		instanceID = config.SearchValue
	} else if config.SearchMode == SearchModeMenu {
		log.Info("menu not implemented yet, must supply a unique tag or instance-id")
		os.Exit(1)
	}

	ins, err := NewAwsInstance(sess, instanceID)
	if err != nil {
		log.Fatalf("could not get instance/session: %s", err)
	}

	log.Debugf("looking up ip of: %s", instanceID)
	sshHost, err := ins.IP(config.ConnectPublicIP)
	if err != nil {
		log.Fatalf("could not find instance ip address: %s", err)
	}

	sshKeyPair, err := NewSSHKeyPair(2048)
	if err != nil {
		log.Fatal(err)
	}

	sshAgent, err := NewSSHAgent()
	if err != nil {
		log.Fatal(err)
	}

	comment := fmt.Sprintf("%s:%s", config.UserName, instanceID)
	err = sshAgent.addKey(sshKeyPair.private, comment)
	if err != nil {
		log.Fatal(err)
	}

	sshArgs = append(sshArgs, sshHost)
	sshArgs = append(sshArgs, sshExtraArgs...)

	log.Debugf("host: %s", sshHost)

	log.Debug("pushing public key to instance")
	err = ins.sendPublicKey(config.UserName, string(sshKeyPair.public))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("running command: ssh %s", strings.Join(sshArgs[:], " "))
	err = sshConnect(sshArgs)
	if err != nil {
		log.Fatal(err)
	}
}
