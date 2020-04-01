package main

import (
	"errors"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

type SearchMode int

const (
	SearchModeMenu SearchMode = 1 << iota
	SearchModeInst
	SearchModeTag
)

type ESSHConfig struct {
	UserName        string
	Region          string
	ConnectPublicIP bool
	Debug           bool
	SearchMode      SearchMode

	// Search value will either be a instance id or tag, check SearchMode to find out what
	SearchValue  string
	sshExtraArgs []string
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
	}

	if config.Region == "" {
		return nil, errors.New("count not find your AWS region from either -r or env vars AWS_REGION, AWS_DEFAULT_REGION")
	}

	config.ConnectPublicIP = *usePublicIP

	if config.Debug {
		log.Debug("All cmd line args passed in")
		for idx := range flag.Args() {
			log.Debugf("flag_pos: %d, flag: %s\n", idx, flag.Arg(idx))
		}
	}

	nargs := flag.NArg()

	// Do we have extra flags denoted by a --
	lastDashAt := flag.CommandLine.ArgsLenAtDash()

	if lastDashAt != -1 && lastDashAt > 1 {
		return nil, errors.New("only specifiy an instance id or a tag, if a tag has a space, wrap in double quotes")
	}

	if lastDashAt == 0 || lastDashAt == 1 {
		config.sshExtraArgs = flag.Args()[lastDashAt:flag.NArg()]
	}

	if lastDashAt > 0 || (lastDashAt == -1 && nargs == 1) {
		config.SearchValue = flag.Arg(0)
		if strings.HasPrefix(config.SearchValue, "i-") {
			config.SearchMode = SearchModeInst
		} else {
			config.SearchMode = SearchModeTag
		}
	} else if nargs > 1 && lastDashAt != 0 {
		return nil, errors.New("only specifiy an instance id or a tag, if a tag has a space, wrap in double quotes")
	} else {
		config.SearchMode = SearchModeMenu
	}

	return config, nil
}
