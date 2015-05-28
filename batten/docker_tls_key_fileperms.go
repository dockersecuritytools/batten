package batten

import "errors"

func (dc *DockerTLSKeyFilePermsCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerTLSKeyFilePermsCheck) AuditCheck() (bool, error) {
	succ, args, err := readDockerDaemonArgs(dc.dockerPidFile)

	if err != nil {
		return false, err
	}

	if succ {
		return dc.validateFromArgs("--tlskey", args)
	}

	return false, errors.New("Docker daemon not running")
}

type DockerTLSKeyFilePermsCheck struct {
	*CheckDefinitionImpl
	*FilePermsCheck
	dockerPidFile string
}

func makeDockerTLSKeyFilePermsCheck() Check {
	return &DockerTLSKeyFilePermsCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-3.24",
			category:    `Docker daemon configuration files`,
			name:        `Verify that Docker server certificate key file permissions are set to 400`,
			description: `Verify that the Docker server certificate key file (the file that is passed alongwith '--tlskey' parameter) has permissions of '400'.`,
			rationale:   `The Docker server certificate key file should be protected from any tampering or unneeded reads. It holds the private key for the Docker server certificate. Hence, it must have permissions of '400' to maintain the integrity of the Docker server certificate.`,
			auditDescription: `Execute the below command to verify that the Docker server certificate key file has permissions of '400':

stat -c %a <path to Docker server certificate key file>`,
			remediation: `#> chmod 400 <path to Docker server certificate key file>

This would set the Docker server certificate key file permissions to '400'.`,
			impact:       `None.`,
			defaultValue: `By default, the permissions for Docker server certificate key file might not be '400'. The default file permissions are governed by the system or user specific umask values. `,
			references: []string{
				"https://docs.docker.com/articles/certificates/",
				"http://docs.docker.com/articles/https/",
			},
		},
		FilePermsCheck: &FilePermsCheck{targetPerms: 0400},
	}
}
