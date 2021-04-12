package main

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/z0mbix/essh/internal/aws"
	"github.com/z0mbix/essh/internal/config"
	"github.com/z0mbix/essh/internal/menu"
	"github.com/z0mbix/essh/internal/ssh"
)

var (
	version = "dev"
)

func main() {
	if config.ShowVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	sess, err := aws.NewSession(config.Region)
	if err != nil {
		log.Fatalf("could not get instance/session: %s", err)
	}

	instance, err := menu.GetInstance(sess)
	if err != nil {
		log.Fatalf("%s", err)
	}
	log.Debugf("host: %s", instance.ConnectIP)

	comment := fmt.Sprintf("%s:%s", config.UserName, instance.ID)
	essh, err := ssh.NewSession(comment, config.PrivateKeyLifetime)
	if err != nil {
		log.Fatalf("%s", err)
	}

	sshArgs := []string{"-l", config.UserName}
	sshArgs = append(sshArgs, instance.ConnectIP)
	sshArgs = append(sshArgs, config.ExtraArgs...)

	err = essh.Connect(instance, sshArgs)
	if err != nil {
		log.Fatal(err.Error())
	}
}
