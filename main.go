package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

func main() {
	userName := flag.StringP("username", "u", "ec2-user", "UNIX user name")
	region := flag.StringP("region", "r", "", "AWS Region")
	usePublicIP := flag.BoolP("use-public-ip", "p", false, "Use the public ip instead of the private ip address")
	debug := flag.BoolP("debug", "d", false, "Enable debug logging")
	flag.Parse()

	log.SetLevel(log.InfoLevel)
	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000",
		FullTimestamp:   true,
	})

	var instanceID string
	var hasInstanceID bool

	awsRegion := *region
	if awsRegion == "" {
		log.Debug("aws region not set, trying AWS_DEFAULT_REGION environment variable")
		awsRegion = os.Getenv("AWS_DEFAULT_REGION")
		if awsRegion == "" {
			log.Debug("aws region not found in AWS_DEFAULT_REGION environment variable")
			log.Fatal("please set the region using the -r/--region flag or the AWS_DEFAULT_REGION environment variable")
		}
		log.Debugf("aws region found in AWS_DEFAULT_REGION environment variable: %s", awsRegion)
	}

	sshHost := flag.Arg(0)
	sshArgs := []string{"-l", *userName}
	sshExtraArgs := flag.Args()[1:len(flag.Args())]

	if strings.HasPrefix(sshHost, "i-") {
		hasInstanceID = true
		instanceID = sshHost
	}

	if !hasInstanceID {
		log.Fatal("name lookup not yet supported")
	}

	ins, err := NewAwsInstance(awsRegion, instanceID)
	if err != nil {
		log.Fatalf("could not get instance/session: %s", err)
	}

	log.Debugf("looking up ip of instance id: %s", sshHost)
	sshHost, err = ins.IP(*usePublicIP)
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

	comment := fmt.Sprintf("%s:%s", *userName, instanceID)
	err = sshAgent.addKey(sshKeyPair.private, comment)
	if err != nil {
		log.Fatal(err)
	}

	sshArgs = append(sshArgs, sshHost)
	sshArgs = append(sshArgs, sshExtraArgs...)

	log.Debugf("host: %s", sshHost)

	log.Debug("pushing public key to instance")
	err = ins.sendPublicKey(*userName, string(sshKeyPair.public))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("running command: ssh %s", strings.Join(sshArgs[:], " "))
	err = sshConnect(sshArgs)
	if err != nil {
		log.Fatal(err)
	}
}
