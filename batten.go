package main

import (
	"os"
	"fmt"
	"github.com/dockersecuritytools/batten/batten"
	"github.com/dockersecuritytools/batten/cli"
	"github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v1"
)

const (
	Name        = "batten"
	Description = "Hardening and Auditing Tool For Docker Hosts & Containers"
	Version     = "0.1.0"
)

var (
	app      = kingpin.New(Name, Description)
	// appDebug = app.Flag("debug", "Enable debug mode.").Bool()
	serverIP  = app.Flag("server", "Connect to remote host.").String()
	tlscacert = app.Flag("tlscacert", "TLS CA Certificate.").String()
	tlscert   = app.Flag("tlscert", "TLS Certificate.").String()
	tlskey    = app.Flag("tlskey", "TLS Key.").String()

	appCheck = app.Command("check", "Check host for known issues.")
)

func fatalf(format string, args ...interface{}) {
	fmt.Printf("* fatal: "+format+"\n", args...)
	os.Exit(1)
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stderr)
}


func main() {
	kingpin.Version(Version)
	args, err := app.Parse(os.Args[1:])
	
	switch kingpin.MustParse(args, err) {
	case appCheck.FullCommand():
		if len(*serverIP) > 0 {
			remoteCheck()
		} else {
			for i, check := range batten.Checks {
				results := batten.RunCheck(check)
				cli.FormatResultsForConsole(i, results)
			}
		}
	default:
		app.Usage(os.Stdout)
	}
}
