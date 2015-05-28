package batten

import (
	"errors"
	"strings"
)

func (dc *DockerEnableIptablesCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerEnableIptablesCheck) lookForIptables(argv string) (bool, error) {
	if strings.Contains(argv, "--iptables=false") {
		return true, nil
	}
	if strings.Contains(argv, "--iptables false") {
		return true, nil
	}
	if strings.Contains(argv, "--iptables=") {
		return false, nil
	}

	return true, nil
}

func (dc *DockerEnableIptablesCheck) AuditCheck() (bool, error) {
	succ, args, err := readDockerDaemonArgs(dc.dockerPidFile)

	if err != nil {
		return false, err
	}

	if succ {
		argv := strings.Join(args, " ")
		return dc.lookForIptables(argv)
	}
	return false, errors.New("Docker daemon not running")
}

type DockerEnableIptablesCheck struct {
	*CheckDefinitionImpl
	dockerPidFile string
}

func makeDockerEnableIptablesCheck() Check {
	return &DockerEnableIptablesCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-2.4",
			category:    `Docker daemon configuration`,
			name:        `Allow Docker to make changes to iptables`,
			description: `Iptables are used to set up, maintain, and inspect the tables of IP packet filter rules in the Linux kernel. Allow the Docker daemon to make changes to the iptables.`,
			rationale:   `Docker will never make changes to your system iptables rules if you choose to do so. Docker server would automatically make the needed changes to iptables based on how you choose your networking options for the containers if it is allowed to do so. It is recommended to let Docker server make changes to iptables automatically to avoid networking misconfiguration that might hamper the communication between containers and to the outside world. Additionally, it would save you hassles of updating iptables every time you choose to run the containers or modify networking options.`,
			auditDescription: `$> ps -ef | grep docker
			Ensure that the '--iptables' parameter is either not present or not set to 'false'.`,
			remediation: `Do not run the Docker daemon with '--iptables=false' parameter. For example, do not start the Docker daemon as below:
			￼
$> docker -d --iptables=false`,
			impact:       `None`,
			defaultValue: `By default, 'iptables' is set to 'true'.`,
			references: []string{
				"http://docs.docker.com/articles/networking/#communication-between-containers",
			},
		},
		dockerPidFile: "/var/run/docker.pid",
	}
}
