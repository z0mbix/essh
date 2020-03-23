package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
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

	var instanceID string
	var hasInstanceID bool

	spew.Dump(config)
	os.Exit(1)

	// awsRegion := *region
	// if awsRegion == "" {
	// 	log.Debug("aws region not set, trying AWS_DEFAULT_REGION environment variable")
	// 	awsRegion = os.Getenv("AWS_DEFAULT_REGION")
	// 	if awsRegion == "" {
	// 		log.Debug("aws region not found in AWS_DEFAULT_REGION environment variable")
	// 		log.Fatal("please set the region using the -r/--region flag or the AWS_DEFAULT_REGION environment variable")
	// 	}
	// 	log.Debugf("aws region found in AWS_DEFAULT_REGION environment variable: %s", awsRegion)
	// }

	if len(flag.Args()) < 1 {
		flag.Usage()
		log.Fatal("You need to specify either an instance id, or a EC2 tag Name")
	}

	sshHost := flag.Arg(0)
	sshArgs := []string{"-l", config.UserName}
	sshExtraArgs := flag.Args()[1:len(flag.Args())]

	if strings.HasPrefix(sshHost, "i-") {
		hasInstanceID = true
		instanceID = sshHost
	}

	sess, err := NewAwsSession(config.Region)
	if err != nil {
		log.Fatalf("could not get instance/session: %s", err)
	}

	if !hasInstanceID {
		log.Debugf("using Name tag %s to find instance id", sshHost)
		instanceID, err = getInstanceIDFromNameTag(sess, sshHost)
		if err != nil {
			log.Fatal(err)
		}
		log.Debugf("found instance id: %s", instanceID)
	}

	ins, err := NewAwsInstance(sess, instanceID)
	if err != nil {
		log.Fatalf("could not get instance/session: %s", err)
	}

	log.Debugf("looking up ip of: %s", sshHost)
	sshHost, err = ins.IP(config.ConnectPublicIP)
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
