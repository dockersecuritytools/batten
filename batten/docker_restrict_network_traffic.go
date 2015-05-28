package batten

import "strings"

func (dc *DockerRestrictedNetworkTrafficCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerRestrictedNetworkTrafficCheck) lookForIccFlag(argv string) (bool, error) {
	if strings.Contains(argv, "--icc=false") {
		return true, nil
	}
	return false, nil
}

func (dc *DockerRestrictedNetworkTrafficCheck) AuditCheck() (bool, error) {
	succ, args, err := readDockerDaemonArgs(dc.dockerPidFile)

	if err != nil {
		return false, err
	}

	if succ {
		argv := strings.Join(args, " ")
		return dc.lookForIccFlag(argv)
	}
	return false, nil
}

type DockerRestrictedNetworkTrafficCheck struct {
	*CheckDefinitionImpl
	dockerPidFile string
}

func makeDockerRestrictedNetworkTrafficCheck() Check {
	return &DockerRestrictedNetworkTrafficCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:   "CIS-Docker-Benchmark-2.2",
			category:     "Docker Daemon Configuration",
			name:         "Restrict network traffic between containers",
			impact:       "The inter container communication would be disabled. No containers would be able to talk to another container on the same host. If any communication between containers on the same host is desired, then it needs to be explicitly defined using container linking.",
			description:  "By default, all network traffic is allowed between containers on the same host. If not desired, restrict all the inter container communication. Link specific containers together that require inter communication.",
			rationale:    "By default, unrestricted network traffic is enabled between all containers on the same host. Thus, each container has the potential of reading all packets across the container network on the same host. This might lead to unintended and unwanted disclosure of information to other containers. Hence, restrict the inter container communication.",
			defaultValue: `By default, all inter container communication is allowed.`,
			auditDescription: `$> ps -ef | grep docker

Ensure that the '--icc' parameter is set to 'false'.`,
			references: []string{
				"https://docs.docker.com/articles/networking",
			},
			remediation: `Run the docker in daemon mode and pass '--icc=false' as argument. For Example,

$> docker -d --icc=false`,
		},
		dockerPidFile: "/var/run/docker.pid",
	}
}
