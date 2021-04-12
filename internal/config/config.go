package config

import (
	"errors"
	"os"
	"strings"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	flag "github.com/spf13/pflag"
)

// SearchModeType is the type of instance search to perform
type SearchModeType int

const (
	SearchModeMenu SearchModeType = 1 << iota
	SearchModeInst
	SearchModeTag
)

var (
	UserName           string
	Region             string
	ConnectPublicIP    bool
	Debug              bool
	SearchMode         SearchModeType
	ShowVersion        bool
	PrivateKeyLifetime uint32

	// SearchValue will either be an instance id or a tag. Check SearchMode to find out which
	SearchValue string
	ExtraArgs   []string

	InvalidRegionErr = errors.New("could not find your AWS region from either -r or environment variables AWS_REGION, AWS_DEFAULT_REGION")
	InvalidArgsErr   = errors.New("only specify an instance id or a tag, if a tag has a space, wrap it in double quotes")
)

func init() {
	region := flag.StringP("region", "r", "", "AWS Region")
	flag.StringVarP(&UserName, "username", "u", "ec2-user", "UNIX user name")
	flag.BoolVarP(&Debug, "debug", "d", false, "Enable debug logging")
	flag.BoolVarP(&ConnectPublicIP, "use-public-ip", "p", false, "Use the public ip instead of the private ip address")
	flag.Uint32VarP(&PrivateKeyLifetime, "key-ttl", "t", 10, "How long the private key will live in the ssh-agent in seconds")
	flag.BoolVarP(&ShowVersion, "version", "v", false, "Show version")

	// Remove the annoying "pflag: help requested" message
	flag.ErrHelp = errors.New("")
	flag.Parse()

	log.SetHandler(cli.New(os.Stdout))

	if Debug {
		log.SetLevel(log.DebugLevel)
	}

	Region = configureRegion(*region)
	if Region == "" {
		log.Fatalf("%s", InvalidRegionErr)
	}

	if user := os.Getenv("ESSH_DEFAULT_USER"); user != "" {
		log.Debugf("Setting region from ESSH_DEFAULT_USER env: %s", user)
		UserName = user
	}

	if Debug {
		log.Debug("all cmd line args passed in:")
		for idx := range flag.Args() {
			log.Debugf("flag position: %d, flag: %s", idx, flag.Arg(idx))
		}
	}

	nargs := flag.NArg()

	// Do we have extra flags denoted by a --
	lastDashAt := flag.CommandLine.ArgsLenAtDash()

	if lastDashAt != -1 && lastDashAt > 1 {
		log.Fatalf("%s", InvalidArgsErr)
	}

	if lastDashAt == 0 || lastDashAt == 1 {
		ExtraArgs = flag.Args()[lastDashAt:flag.NArg()]
	}

	if lastDashAt > 0 || (lastDashAt == -1 && nargs == 1) {
		SearchValue = flag.Arg(0)
		if strings.HasPrefix(SearchValue, "i-") {
			SearchMode = SearchModeInst
		} else {
			SearchMode = SearchModeTag
		}
	} else if nargs > 1 && lastDashAt != 0 {
		log.Fatalf("%s", InvalidArgsErr)
	} else {
		SearchMode = SearchModeMenu
	}
}

func configureRegion(arg string) string {
	if arg != "" {
		log.Debugf("setting region from args: %s", arg)
		return arg
	}

	if region := os.Getenv("AWS_REGION"); region != "" {
		log.Debugf("Setting region from AWS_REGION env: %s", region)
		return region
	}

	if region := os.Getenv("AWS_DEFAULT_REGION"); region != "" {
		log.Debugf("Setting region from AWS_DEFAULT_REGION env: %s", region)
		return region
	}

	return ""
}
