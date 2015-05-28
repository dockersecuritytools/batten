package batten

import "errors"

func (dc *DockerTLSCACertOwnerCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerTLSCACertOwnerCheck) validate(args []string) (bool, error) {
	lookFor := "--tlscacert"
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

func (dc *DockerTLSCACertOwnerCheck) AuditCheck() (bool, error) {
	succ, args, err := readDockerDaemonArgs(dc.dockerPidFile)

	if err != nil {
		return false, err
	}

	if succ {
		return dc.validate(args)
	}

	return false, errors.New("Docker daemon not running")
}

type DockerTLSCACertOwnerCheck struct {
	*CheckDefinitionImpl
	*FileOwnerCheck
	dockerPidFile string
}

func makeDockerTLSCACertOwnerCheck() Check {
	return &DockerTLSCACertOwnerCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-3.19",
			category:    `Docker daemon configuration files`,
			name:        `Verify that TLS CA certificate file ownership is set to root:root`,
			description: `Verify that the TLS CA certificate file (the file that is passed alongwith '--tlscacert' parameter) is owned and group-owned by 'root'.`,
			rationale:   `The TLS CA certificate file should be protected from any tampering. It is used to authenticate Docker server based on given CA certificate. Hence, it must be owned and group-owned by 'root' to maintain the integrity of the CA certificate.`,
			auditDescription: `Execute the below command to verify that the TLS CA certificate file is owned and group- owned by 'root':
stat -c %U:%G <path to TLS CA certificate file> | grep -v root:root
			￼
The above command should not return anything.`,
			remediation: `#> chown root:root <path to TLS CA certificate file>

This would set the ownership and group-ownership for the TLS CA certificate file to 'root'.`,
			impact:       `None`,
			defaultValue: `By default, the ownership and group-ownership for TLS CA certificate file is correctly set to 'root'.`,
			references: []string{
				"https://docs.docker.com/articles/certificates/",
				"http://docs.docker.com/articles/https/",
			},
		},
		FileOwnerCheck: &FileOwnerCheck{uid: 0, gid: 0},
	}
}
