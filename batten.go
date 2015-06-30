package main

import (
	"os"

	"github.com/dockersecuritytools/batten/batten"
	"github.com/dockersecuritytools/batten/cli"
	"gopkg.in/alecthomas/kingpin.v1"
)

const (
	Name        = "batten"
	Description = "Hardening and Auditing Tool For Docker Hosts & Containers"
	Version     = "0.1.0"
)

var (
	app      = kingpin.New(Name, Description)
	appDebug = kingpin.Flag("debug", "Enable debug mode.").Bool()

	appCheck = app.Command("check", "Check host for known issues.")
)

func main() {
	kingpin.Version(Version)
	args, err := app.Parse(os.Args[1:])

	switch kingpin.MustParse(args, err) {
	case appCheck.FullCommand():
		for i, check := range batten.Checks {
			results := batten.RunCheck(check)

			cli.FormatResultsForConsole(i, results)
		}
	default:
		app.Usage(os.Stdout)
	}

}
