package config

import (
	"errors"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
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
	UserName        string
	Region          string
	ConnectPublicIP bool
	Debug           bool
	SearchMode      SearchModeType
	ShowVersion     bool

	// SearchValue will either be a instance id or tag, check SearchMode to find out what
	SearchValue string
	ExtraArgs   []string

	InvalidRegionErr = errors.New("count not find your AWS region from either -r or env vars AWS_REGION, AWS_DEFAULT_REGION")
	InvalidArgsErr   = errors.New("only specify an instance id or a tag, if a tag has a space, wrap in double quotes")
)

func init() {
	region := flag.StringP("region", "r", "", "AWS Region")

	flag.StringVarP(&UserName, "username", "u", "ec2-user", "UNIX user name")
	flag.BoolVarP(&Debug, "debug", "d", false, "Enable debug logging")
	flag.BoolVarP(&ConnectPublicIP, "use-public-ip", "p", false, "Use the public ip instead of the private ip address")
	flag.BoolVarP(&ShowVersion, "version", "v", false, "Show version")

	flag.Parse()

	log.SetLevel(log.InfoLevel)

	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp:       true,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})
	if Debug {
		log.SetLevel(log.DebugLevel)
	}

	Region = configureRegion(*region)

	if Debug {
		log.Debug("All cmd line args passed in")
		for idx := range flag.Args() {
			log.Debugf("flag_pos: %d, flag: %s\n", idx, flag.Arg(idx))
		}
	}

	nargs := flag.NArg()

	// Do we have extra flags denoted by a --
	lastDashAt := flag.CommandLine.ArgsLenAtDash()

	if lastDashAt != -1 && lastDashAt > 1 {
		log.Fatal(InvalidArgsErr)
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
		log.Fatal(InvalidArgsErr)
	} else {
		SearchMode = SearchModeMenu
	}
	log.Debug(spew.Sdump())
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

	return "eu-west-1"
}
