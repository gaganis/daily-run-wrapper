package main

import (
	"cloud.google.com/go/civil"
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"
)

type Configuration struct {
	Profile             string
	Command             string
	Interval            time.Duration
	HasPreferredRunTime bool
	PreferredRunTime    civil.Time
	RunNow              bool
}

func ParseConfigFromFlags() Configuration {
	configuration := Configuration{}
	flag.StringVar(&configuration.Profile, "profile", "default", "Profile to use.")

	flag.StringVar(&configuration.Command, "command",
		"echo 'daily-runner has run echo printing this text'",
		"The command that runner will execute")

	flag.DurationVar(&configuration.Interval, "interval",
		4*time.Minute,
		"The interval that daily-runner will use to check if it needs to run. "+
			"Can accept values acceptable to golang time.ParseDuration function")

	flag.BoolVar(&configuration.RunNow, "now",
		false,
		"Run the command immediately. Previous runs are not taken into account. "+
			"This option also disables process locking, so multiple processes can run at the same time.")

	preferedRunTimePtr := flag.String("preferredTime", "",
		"Set a preferred time for the runner to run command. This time overrides the daily logic and the command will always "+
			"run if the system is up at that time.")

	flag.Parse()
	isValid, message := validateProfile(configuration.Profile)
	if !isValid {
		fmt.Print(message)
		os.Exit(2)
	}

	if *preferedRunTimePtr != "" {
		parsedCivilTime, err2 := civil.ParseTime(*preferedRunTimePtr)
		if err2 != nil {
			fmt.Printf("Parsing PreferredRunTime flag failed with error: %v", err2)
			os.Exit(2)
		}
		configuration.HasPreferredRunTime = true
		configuration.PreferredRunTime = parsedCivilTime
	} else {
		configuration.HasPreferredRunTime = false
	}
	return configuration
}

func validateProfile(profile string) (bool, string) {
	matched, _ := regexp.MatchString(`^[0-9a-zA-Z.\-_]*$`, profile)
	if !matched {
		return false, fmt.Sprintf("Wrong profile name argument provided: '%v'. Please provide a name containing only latin "+
			"letters, numbers and the characters '.-_'. Exiting", profile)
	}

	return true, ""
}
