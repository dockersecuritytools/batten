package batten

import "errors"

func (dc *DockerTLSCertFilePermsCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerTLSCertFilePermsCheck) AuditCheck() (bool, error) {
	succ, args, err := readDockerDaemonArgs(dc.dockerPidFile)

	if err != nil {
		return false, err
	}

	if succ {
		return dc.validateFromArgs("--tlscert", args)
	}

	return false, errors.New("Docker daemon not running")
}

type DockerTLSCertFilePermsCheck struct {
	*CheckDefinitionImpl
	*FilePermsCheck
	dockerPidFile string
}

func makeDockerTLSCertFilePermsCheck() Check {
	return &DockerTLSCertFilePermsCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-3.22",
			category:    `Docker daemon configuration files`,
			name:        `Verify that Docker server certificate file permissions are set to 444 or more restrictive`,
			description: `Verify that the Docker server certificate file (the file that is passed alongwith '-- tlscert' parameter) has permissions of '444' or more restrictive.`,
			rationale:   `The Docker server certificate file should be protected from any tampering. It is used to authenticate Docker server based on the given server certificate. Hence, it must be have permissions of '444' to maintain the integrity of the certificate.`,
			auditDescription: `Execute the below command to verify that the Docker server certificate file has permissions of '444' or more restrictive:

stat -c %a <path to Docker server certificate file>`,
			remediation: `#> chmod 444 <path to Docker server certificate file>

This would set the file permissions of the Docker server file to '444'.`,
			impact:       `None`,
			defaultValue: `By default, the permissions for Docker server certificate file might not be '444'. The default file permissions are governed by the system or user specific umask values.`,
			references: []string{
				"https://docs.docker.com/articles/certificates/",
				"http://docs.docker.com/articles/https/",
			},
		},
		FilePermsCheck: &FilePermsCheck{targetPerms: 0444},
	}
}
