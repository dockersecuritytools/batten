package batten

import "errors"

func (dc *DockerTLSCertOwnerCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerTLSCertOwnerCheck) validate(args []string) (bool, error) {
	lookFor := "--tlscert"
	filepath := getArgValue(lookFor, args)

	if filepath != "" {
		if PathExists(filepath) {
			dc.filepath = filepath
			return dc.IsOwnerAndGroupOwner(dc.uid, dc.gid)
		} else {
			return false, errors.New("Could not find path: " + filepath)
		}
	}
	// it's ok, no tls config set
	return true, nil
}

func (dc *DockerTLSCertOwnerCheck) AuditCheck() (bool, error) {
	succ, args, err := readDockerDaemonArgs(dc.dockerPidFile)

	if err != nil {
		return false, err
	}

	if succ {
		return dc.validate(args)
	}

	return false, errors.New("Docker daemon not running")
}

type DockerTLSCertOwnerCheck struct {
	*CheckDefinitionImpl
	*FileOwnerCheck
	dockerPidFile string
}

func makeDockerTLSCertOwnerCheck() Check {
	return &DockerTLSCertOwnerCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-3.21",
			category:    `Docker daemon configuration files`,
			name:        `Verify that Docker server certificate file ownership is set to root:root`,
			description: `Verify that the Docker server certificate file (the file that is passed alongwith '-- tlscert' parameter) is owned and group-owned by 'root'.`,
			rationale:   `The Docker server certificate file should be protected from any tampering. It is used to authenticate Docker server based on the given server certificate. Hence, it must be owned and group-owned by 'root' to maintain the integrity of the certificate.`,
			auditDescription: `Execute the below command to verify that the Docker server certificate file is owned and group-owned by 'root':
			￼
stat -c %U:%G <path to Docker server certificate file> | grep -v root:root

The above command should not return anything.`,
			remediation: `#> chown root:root <path to Docker server certificate file>

This would set the ownership and group-ownership for the Docker server certificate file to 'root'.`,
			impact:       `None`,
			defaultValue: `By default, the ownership and group-ownership for Docker server certificate file is correctly set to 'root'.`,
			references: []string{
				"https://docs.docker.com/articles/certificates/",
				"http://docs.docker.com/articles/https/",
			},
		},
		FileOwnerCheck: &FileOwnerCheck{},
	}

}
