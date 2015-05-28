package batten

import "errors"

func (dc *DockerTLSKeyOwnerCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerTLSKeyOwnerCheck) validate(args []string) (bool, error) {
	lookFor := "--tlskey"
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

func (dc *DockerTLSKeyOwnerCheck) AuditCheck() (bool, error) {

	succ, args, err := readDockerDaemonArgs(dc.dockerPidFile)

	if err != nil {
		return false, err
	}

	if succ {
		return dc.validate(args)
	}

	return false, errors.New("Docker daemon not running")
}

type DockerTLSKeyOwnerCheck struct {
	*CheckDefinitionImpl
	*FileOwnerCheck
	dockerPidFile string
}

func makeDockerTLSKeyOwnerCheck() Check {
	return &DockerTLSKeyOwnerCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-3.23",
			category:    `Docker daemon configuration files`,
			name:        `Verify that Docker server certificate key file ownership is set to root:root`,
			description: `Verify that the Docker server certificate key file (the file that is passed alongwith '--tlskey' parameter) is owned and group-owned by 'root'.`,
			rationale:   `The Docker server certificate key file should be protected from any tampering or unneeded reads. It holds the private key for the Docker server certificate. Hence, it must be owned and group-owned by 'root' to maintain the integrity of the Docker server certificate.`,
			auditDescription: `Execute the below command to verify that the Docker server certificate key file is owned and group-owned by 'root':

stat -c %U:%G <path to Docker server certificate key file> | grep -v root:root

The above command should not return anything.`,
			remediation: `#> chown root:root <path to Docker server certificate key file>

This would set the ownership and group-ownership for the Docker server certificate key file to 'root'.`,
			impact:       `None`,
			defaultValue: `By default, the ownership and group-ownership for Docker server certificate key file is correctly set to 'root'.`,
			references: []string{
				"https://docs.docker.com/articles/certificates/",
				"http://docs.docker.com/articles/https/",
			},
		},
		FileOwnerCheck: &FileOwnerCheck{uid: 0, gid: 0},
	}
}
