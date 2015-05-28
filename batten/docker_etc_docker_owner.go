package batten

func (dc *DockerEtcDockerOwnerCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerEtcDockerOwnerCheck) AuditCheck() (bool, error) {
	if PathExists(dc.filepath) {
		return dc.IsOwnerAndGroupOwner(0, 0)
	}
	return true, nil
}

type DockerEtcDockerOwnerCheck struct {
	*CheckDefinitionImpl
	*FileOwnerCheck
}

func makeDockerEtcDockerOwnerCheck() Check {
	return &DockerEtcDockerOwnerCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-3.15",
			category:    `Docker daemon configuration files`,
			name:        `Verify that /etc/docker directory ownership is set to root:root`,
			description: `Verify that the /etc/docker directory ownership and group-ownership is correctly set to 'root'.`,
			rationale:   `'/etc/docker' directory contains certificates and keys in addition to various sensitive files. Hence, it should be owned and group-owned by 'root' to maintain the integrity of the directory.`,
			auditDescription: `Execute the below command to verify that the directory is owned and group-owned by 'root':

stat -c %U:%G /etc/docker | grep -v root:root

The above command should not return anything. `,
			remediation: `#> chown root:root /etc/docker

This would set the ownership and group-ownership for the directory to 'root'.`,
			impact:       `None`,
			defaultValue: `By default, the ownership and group-ownership for this directory is correctly set to 'root'.`,
			references: []string{
				"https://docs.docker.com/articles/certificates/",
			},
		},

		FileOwnerCheck: &FileOwnerCheck{
			filepath: "/etc/docker",
		},
	}
}
