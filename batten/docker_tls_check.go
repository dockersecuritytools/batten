package batten

import "errors"

func (dc *DockerTLSCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

// TODO: could do this more accurately with lsof -i
func (dc *DockerTLSCheck) lookForListeningConfig(args []string) bool {
	var prev string
	for _, arg := range args {
		if len(arg) > 0 {
			if prev == "-H" || prev == "-H=" {
				return true
			}
			prev = arg
		}
	}
	return false
}

func (dc *DockerTLSCheck) lookForTLSConfigs(args []string) bool {

	if !stringInSlice("--tlsverify", args) {
		return false
	}
	// TODO: this is probably actually optional if the cert
	// is signed by a known good ca
	// if !(stringInSlice(args, "--tlscacert")) {
	// return false
	// }
	if !stringInSlice("--tlscert", args) {
		return false
	}
	if !stringInSlice("--tlskey", args) {
		return false
	}
	return true

}

func (dc *DockerTLSCheck) AuditCheck() (bool, error) {

	succ, args, err := readDockerDaemonArgs(dc.dockerPidFile)

	// TODO: also try a lsof -i -p <pid of docker> -a check??

	if err != nil {
		return false, err
	}

	if succ {

		hasListener := dc.lookForListeningConfig(args)

		if hasListener {
			return dc.lookForTLSConfigs(args), nil
		} else {
			// it's ok
			return true, nil
		}
	}

	return false, errors.New("Docker daemon not running")
}

type DockerTLSCheck struct {
	*CheckDefinitionImpl
	dockerPidFile string
}

func makeDockerTLSCheck() Check {
	return &DockerTLSCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-2.9",
			category:    `Docker daemon configuration`,
			name:        `Configure TLS authentication for Docker daemon`,
			description: `It is possible to make the Docker daemon to listen on a specific IP and port and any other Unix socket other than default Unix socket. Configure TLS authentication to restrict access to Docker daemon via IP and Port.`,
			rationale: `By default, Docker daemon binds to a non-networked Unix socket and runs with 'root' privileges. If you change the default docker daemon binding to a TCP port or any other Unix socket, anyone with access to that port or socket can have full access to Docker daemon and in turn to the host system. Hence, you should not bind the Docker daemon to another IP/Port or a Unix socket.
			If you must expose the Docker daemon via a network socket, configure TLS authentication for the daemon and Docker Swarm APIs (if using). This would restrict the connections to your Docker daemon over the network to a limited number of clients who could successfully authenticate over TLS.`,
			auditDescription: `$> ps -ef | grep docker

			Ensure that the below parameters are present:

			• '--tlsverify'
			• '--tlscacert'
			• '--tlscert'
			• '--tlskey'`,
			remediation:  `Follow the steps mentioned in the Docker documentation or other references.`,
			impact:       `You would need to manage and guard certificates and keys for Docker daemon and Docker clients.`,
			defaultValue: `By default, TLS authentication is not configured.`,
			references: []string{
				"http://docs.docker.com/articles/https/",
				"http://www.hnwatcher.com/r/1644394/Intro-to-Docker-Swarm-Part-2-Comfiguration-Modes-and-Requirements",
				"http://www.blackfinsecurity.com/docker-swarm-with-tls-authentication/",
			},
		},
	}
}
