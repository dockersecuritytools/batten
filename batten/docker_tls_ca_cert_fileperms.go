package batten

import "errors"

func (dc *DockerTLSCACertFilePermsCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerTLSCACertFilePermsCheck) AuditCheck() (bool, error) {
	succ, args, err := readDockerDaemonArgs(dc.dockerPidFile)

	if err != nil {
		return false, err
	}

	if succ {
		return dc.validateFromArgs("--tlscacert", args)
	}

	return false, errors.New("Docker daemon not running")
}

type DockerTLSCACertFilePermsCheck struct {
	*CheckDefinitionImpl
	*FilePermsCheck
	dockerPidFile string
}

func makeDockerTLSCACertFilePermsCheck() Check {
	return &DockerTLSCACertFilePermsCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-3.20",
			category:    `Docker daemon configuration files`,
			name:        `Verify that TLS CA certificate file permissions are set to 444 or more restrictive`,
			description: `Verify that the TLS CA certificate file (the file that is passed along with '--tlscacert' parameter) has permissions of '444' or more restrictive.`,
			rationale:   `The TLS CA certificate file should be protected from any tampering. It is used to authenticate Docker server based on given CA certificate. Hence, it must be have permissions of '444' to maintain the integrity of the CA certificate.`,
			auditDescription: `Execute the below command to verify that the TLS CA certificate file has permissions of '444' or more restrictive:
			￼
stat -c %a <path to TLS CA certificate file>`,
			remediation: `#> chmod 444 <path to TLS CA certificate file>

This would set the file permissions of the TLS CA file to '444'.`,
			defaultValue: `By default, the permissions for TLS CA certificate file might not be '444'. The default file permissions are governed by the system or user specific umask values.`,
			impact:       `None`,
			references: []string{
				"https://docs.docker.com/articles/certificates/",
				"http://docs.docker.com/articles/https/",
			},
		},
		FilePermsCheck: &FilePermsCheck{targetPerms: 0444},
	}
}
