package batten

import (
	"errors"
	"strings"
)

func (dc *DockerPortCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerPortCheck) lookForPorts(args []string) (bool, error) {
	var prev string
	for _, arg := range args {
		if len(arg) > 0 {
			if prev == "-H" {
				if !stringInSlice(arg, dc.whiteListed) {
					// TODO log something?
					return false, nil
				}
			} else if strings.HasPrefix(arg, "-H=") {
				arg = arg[3:]
				if !stringInSlice(arg, dc.whiteListed) {
					// TODO log something?
					return false, nil
				}
			}

			prev = arg
		}
	}
	return true, nil
}

func (dc *DockerPortCheck) AuditCheck() (bool, error) {

	succ, args, err := readDockerDaemonArgs(dc.dockerPidFile)

	// TODO: also try a lsof -i -p <pid of docker> -a check??

	if err != nil {
		return false, err
	}

	if succ {
		return dc.lookForPorts(args)
	}
	return false, errors.New("Docker daemon not running")
}

type DockerPortCheck struct {
	*CheckDefinitionImpl
	// TODO: make configurable
	dockerPidFile string
	whiteListed   []string
}

func makeDockerPortCheck() Check {
	return &DockerPortCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-2.8",
			category:    "Docker Daemon Configuration",
			name:        `Do not bind Docker to another IP/Port or a Unix socket`,
			description: `It is possible to make the Docker daemon to listen on a specific IP and port and any other Unix socket other than default Unix socket. Do not bind Docker daemon to another IP/Port or a Unix socket.`,
			rationale:   `By default, Docker daemon binds to a non-networked Unix socket and runs with 'root' privileges. If you change the default docker daemon binding to a TCP port or any other Unix socket, anyone with access to that port or socket can have full access to Docker daemon and in turn to the host system. Hence, you should not bind the Docker daemon to another IP/Port or a Unix socket.`,
			auditDescription: `$> ps -ef | grep docker

Ensure that the '-H' parameter is not present.`,
			remediation: `Do not bind the Docker daemon to any IP and Port or a non-default Unix socket. For example, do not start the Docker daemon as below:

$> docker -H tcp://10.1.2.3:2375 -H unix:///var/run/example.sock -d`,
			impact:       `No one can have full access to Docker daemon except 'root'. Alternatively, you should configure the TLS authentication for Docker and Docker Swarm APIs if you want to bind the Docker daemon to any other IP and Port.`,
			defaultValue: `By default, Docker daemon binds to a non-networked Unix socket.`,
			references: []string{
				"https://docs.docker.com/articles/basics/#bind-docker-to-another-hostport-or-a- unix-socket",
			},
		},
	}
}
