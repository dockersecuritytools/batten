package batten

import (
	"errors"
	"strings"
)

func (dc *DockerSetLoggingLevelCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerSetLoggingLevelCheck) lookForLoggingLevel(argv string) (bool, error) {
	if strings.Contains(argv, "--log-level=info") {
		return true, nil
	}
	if strings.Contains(argv, "--log-level info") {
		return true, nil
	}

	if strings.Contains(argv, "--log-level=") {
		return false, nil
	}

	return true, nil
}

func (dc *DockerSetLoggingLevelCheck) AuditCheck() (bool, error) {
	succ, args, err := readDockerDaemonArgs(dc.dockerPidFile)

	if err != nil {
		return false, err
	}

	if succ {
		argv := strings.Join(args, " ")
		return dc.lookForLoggingLevel(argv)
	}
	return false, errors.New("Docker daemon not running")
}

type DockerSetLoggingLevelCheck struct {
	*CheckDefinitionImpl
	dockerPidFile string
}

func makeDockerSetLoggingLevelCheck() Check {
	return &DockerSetLoggingLevelCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-2.3",
			category:    `Docker daemon configuration`,
			name:        `Set the logging level`,
			description: `Set Docker daemon log level to 'info'.`,
			rationale:   `Setting up an appropriate log level, configures the Docker daemon to log events that you would want to review later. A base log level of 'info' and above would capture all logs except debug logs. Until and unless required, you should not run Docker daemon at 'debug' log level.`,
			auditDescription: `$> ps -ef | grep docker
Ensure that either the '--log-level' parameter is not present or if present, then it is set to 'info'.`,
			remediation: `Run the Docker daemon as below:

$> docker -d --log-level="info"`,
			impact:       `None`,
			defaultValue: `By default, Docker daemon is set to log level of 'info'.`,
			references: []string{
				"https://docs.docker.com/reference/commandline/cli/#daemon",
			},
		},
		dockerPidFile: "/var/run/docker.pid",
	}
}
