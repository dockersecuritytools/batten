package main

import (
	"bytes"
	"os"
	"io"
	"io/ioutil"
	"fmt"
	"log"
	"./batten" 
	"./cli"
	"gopkg.in/alecthomas/kingpin.v1"
	docker "github.com/fsouza/go-dockerclient"
)

const (
	Name        = "batten"
	Description = "Hardening and Auditing Tool For Docker Hosts & Containers"
	Version     = "0.1.0"
	BattenDockerRepository = "jerbi/batten"
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

var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

func initLogs(
    traceHandle   io.Writer,
    infoHandle    io.Writer,
    warningHandle io.Writer,
    errorHandle   io.Writer) {

    Trace = log.New(traceHandle,
        "TRACE: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Warning = log.New(warningHandle,
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

func newDockerClient() (*docker.Client, error) {
	if len(*tlscert) == 0 {
		return docker.NewClient(*serverIP)
	} else {
		return docker.NewTLSClient(*serverIP, *tlscert, *tlskey, *tlscacert)
	}
}

func remoteCheck() (error){
	
	client, err := newDockerClient()
	
	if err != nil {
		Error.Fatalf("Failed to connect to host '%s'. Error: %s", *serverIP, err)
	}
	
	// Pull Batten image from the Docker repository
	var buf bytes.Buffer
	err = client.PullImage(docker.PullImageOptions{Repository: BattenDockerRepository, OutputStream: &buf}, docker.AuthConfiguration{})
	if err != nil {
		Error.Fatalf("Failed to pull '%s' image on host '%s'. Error: %s", BattenDockerRepository, *serverIP, err)
	}
	
	// Create Batten container
	config := &docker.Config{Image: BattenDockerRepository}
	container, err := client.CreateContainer(docker.CreateContainerOptions{Config: config})
	if err != nil {
		return err
	}
			
	// Start Batten container
	binds := []string{"/var/run/docker.sock:/var/run/docker.sock"}
	hostConfig := &docker.HostConfig{Binds: binds}
	err = client.StartContainer(container.ID, hostConfig)
	if err != nil {
		return err
	}	
	
	// Wait for container scan to finish
	code, err := client.WaitContainer(container.ID)
	
	if err != nil {
		return err
	}
	
	if code != 0 {
		fmt.Printf("container run returned %d", code )	
	}
	
	// Print container logs
	opts := docker.AttachToContainerOptions{
		Container:    container.ID,
		OutputStream: &buf,
		Stdout:       true,
		Stderr:       true,
		Logs:         true,
	}
	err = client.AttachToContainer(opts)
	if err != nil {
		return err
	}
	fmt.Printf("%s", buf.String())
	
	// Cleanup - remove container and image
	err = client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID})
	if err != nil {
		err = client.RemoveImage(BattenDockerRepository)
	}
	
	return err
}

func main() {
	initLogs(ioutil.Discard, os.Stdout, os.Stdout, os.Stdout)
		
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
