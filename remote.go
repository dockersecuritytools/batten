package main

import (
	"bytes"
	"fmt"
	"github.com/mgutz/ansi"
	docker "github.com/fsouza/go-dockerclient"
)

const (
	BattenDockerRepository = "jerbi/batten"
)


func colorPrint(color string, format string, args ...interface{}) {
	fmt.Printf(color + format + ansi.Reset + "\n", args...)
}

func newDockerClient() (*docker.Client, error) {
	if len(*tlscert) == 0 {
		return docker.NewClient(*serverIP)
	} else {
		return docker.NewTLSClient(*serverIP, *tlscert, *tlskey, *tlscacert)
	}
}

func cleanUp(client *docker.Client, container *docker.Container) {
	client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID})
	client.RemoveImage(BattenDockerRepository)
}

func remoteCheck() (error){	
	client, err := newDockerClient()
	
	if err != nil {
		fatalf("Failed to connect to host '%s'. Error: %v", *serverIP,err)
	}
	
	// Pull Batten image from the Docker repository
	var pullBuf bytes.Buffer
	colorPrint(ansi.Green, "Pulling image '%s' on host '%s'...", BattenDockerRepository, *serverIP)
	err = client.PullImage(docker.PullImageOptions{Repository: BattenDockerRepository, OutputStream: &pullBuf}, docker.AuthConfiguration{})
	if err != nil {
		fatalf("Failed to pull '%s' image on host '%s'. Error: %v", BattenDockerRepository, *serverIP, err)
	}
	
	// Create Batten container
	config := &docker.Config{Image: BattenDockerRepository}
	container, err := client.CreateContainer(docker.CreateContainerOptions{Config: config})
	if err != nil {
		client.RemoveImage(BattenDockerRepository)
		fatalf("Failed to create container on host '%s'. Error: %v", *serverIP, err)
	}
			
	// Start Batten container
	binds := []string{"/var/run/docker.sock:/var/run/docker.sock"}
	hostConfig := &docker.HostConfig{Binds: binds}
	err = client.StartContainer(container.ID, hostConfig)
	if err != nil {
		cleanUp(client, container)
		fatalf("Failed to start container on host '%s'. Error: %v", *serverIP, err)
	}	
	
	// Wait for container scan to finish
	colorPrint(ansi.Green, "Running scan on host '%s'...", *serverIP)
	code, err := client.WaitContainer(container.ID)
	
	if err != nil || code != 0{
		cleanUp(client, container)
		fatalf("Container finished with errors. Host: '%s', code: '%d', error: %v", *serverIP, code, err)
	}
	
	var logsBuf bytes.Buffer
	// Print container logs
	opts := docker.AttachToContainerOptions{
		Container:    container.ID,
		OutputStream: &logsBuf,
		Stdout:       true,
		Stderr:       true,
		Logs:         true,
	}
	err = client.AttachToContainer(opts)
	if err != nil {
		cleanUp(client, container)
		fatalf("Failed to attach to contianer output on host '%s'. Error: %v", *serverIP, err)
	}
	fmt.Printf("%s", logsBuf.String())
	
	// Cleanup - remove container and image
	cleanUp(client, container)
	
	return err
}

