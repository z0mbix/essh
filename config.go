package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

type ESSHConfig struct {
	UserName        string
	Region          string
	ConnectPublicIP bool
	Debug           bool
	InstanceID      string
}

func defaultConfig() *ESSHConfig {
	return &ESSHConfig{
		UserName:        "ec2-user",
		Region:          "eu-west-1",
		ConnectPublicIP: false,
	}
}

func getESSHConfig() (*ESSHConfig, error) {
	config := defaultConfig()

	userName := flag.StringP("username", "u", "ec2-user", "UNIX user name")
	region := flag.StringP("region", "r", "", "AWS Region")
	usePublicIP := flag.BoolP("use-public-ip", "p", false, "Use the public ip instead of the private ip address")
	debug := flag.BoolP("debug", "d", false, "Enable debug logging")
	flag.Parse()

	config.Debug = *debug
	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	if *userName != "" {
		config.UserName = *userName
	}

	if *region != "" {
		config.Region = *region
	} else {
		config.Region = os.Getenv("AWS_REGION")
		if config.Region == "" {
			config.Region = os.Getenv("AWS_DEFAULT_REGION")
		}
		//TODO: Add error, no region, cannot contunue
	}

	config.ConnectPublicIP = *usePublicIP

	//TODO: check nargs for pos args

	return config, nil
}
