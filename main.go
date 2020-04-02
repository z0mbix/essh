package main

import (
	"fmt"
	"github.com/z0mbix/essh/internal/ssh"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

func main() {
	var err error

	log.SetLevel(log.InfoLevel)

	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp:       true,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})

	config, err := getESSHConfig()

	if err != nil {
		log.Fatal(err)
	}

	if config.Debug {
		log.Debug(spew.Sdump(config))
	}

	sshArgs := []string{"-l", config.UserName}
	sshExtraArgs := config.sshExtraArgs

	sess, err := NewAwsSession(config.Region)
	if err != nil {
		log.Fatalf("could not get instance/session: %s", err)
	}

	var reservations []*ec2.Reservation

	if config.SearchMode == SearchModeTag {
		log.Debugf("using Name tag %s to find instance id", config.SearchValue)

		reservations, err = getInstanceFromNameTag(sess, config.SearchValue)
		if err != nil {
			log.Fatal(err)
		}

	} else if config.SearchMode == SearchModeInst {
		reservations, err = getInstanceFromID(sess, config.SearchValue)
		if err != nil {
			log.Fatal(err)
		}
	} else if config.SearchMode == SearchModeMenu {
		reservations, err = getInstanceFromNameTag(sess, "*")
		if err != nil {
			log.Fatal(err)
		}
	}

	if len(reservations) == 0 {
		log.Fatal("no instance found, add better logging here")
	}

	var instConnect *AwsInstance
	if len(reservations) == 1 && len(reservations[0].Instances) == 1 {
		instConnect, err = NewAwsInstance(sess, reservations[0].Instances[0], config.ConnectPublicIP)
		if err != nil {
			log.Fatalf("could not get instance/session: %s", err)
		}
	} else { //Menu Choices
		instances := []AwsInstance{}
		for rIdx := range reservations {
			for _, inst := range reservations[rIdx].Instances {
				i, err := NewAwsInstance(sess, inst, config.ConnectPublicIP)
				if err != nil {
					log.Fatalf("could not get instance/session: %s", err)
				}
				instances = append(instances, *i)
			}
		}
		instConnect, err = showMenu(instances)
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("looking up ip of: %s", instConnect.ConnectIP)
	sshHost := instConnect.ConnectIP
	if err != nil {
		log.Fatalf("could not find instance ip address: %s", err)
	}

	comment := fmt.Sprintf("%s:%s", config.UserName, instConnect.ID)
	session, err := ssh.NewSession(comment)
	if err != nil {
		log.Fatal(err)
	}

	sshArgs = append(sshArgs, sshHost)
	sshArgs = append(sshArgs, sshExtraArgs...)

	log.Debugf("host: %s", sshHost)

	log.Debug("pushing public key to instance")
	err = instConnect.sendPublicKey(config.UserName, string(session.KeyPair.PublicKey))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("running command: ssh %s", strings.Join(sshArgs[:], " "))
	err = session.Connect(sshArgs)
	if err != nil {
		log.Fatal(err)
	}
}
