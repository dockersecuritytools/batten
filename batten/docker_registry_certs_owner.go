package batten

func (dc *DockerRegistryCertsOwnerCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerRegistryCertsOwnerCheck) AuditCheck() (bool, error) {
	if PathExists(dc.filepath) {
		return dc.IsOwnerAndGroupOwnerRecursive(0, 0)
	}

	return true, nil
}

type DockerRegistryCertsOwnerCheck struct {
	*CheckDefinitionImpl
	*FileOwnerCheck
}

func makeDockerRegistryCertsOwnerCheck() Check {
	return &DockerRegistryCertsOwnerCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier: "CIS-Docker-Benchmark-3.17",
			category:   `Docker daemon configuration files`,

			name:        `Verify that registry certificate file ownership is set to root:root`,
			description: `Verify that all the registry certificate files (usually found under /etc/docker/certs.d/<registry-name> directory) are owned and group-owned by 'root'.`,
			rationale:   `/etc/docker/certs.d/<registry-name> directory contains Docker registry certificates. These certificate files must be owned and group-owned by 'root' to maintain the integrity of the certificates.`,
			auditDescription: `Execute the below command to verify that the registry certificate files are owned and group-owned by 'root':

stat -c %U:%G /etc/docker/certs.d/* | grep -v root:root

The above command should not return anything.`,
			remediation: `#> chown root:root /etc/docker/certs.d/<registry-name>/*

This would set the ownership and group-ownership for the registry certificate files to 'root'.
			`,
			impact:       `None`,
			defaultValue: `By default, the ownership and group-ownership for registry certificate files is correctly set to 'root'.`,
			references: []string{
				"https://docs.docker.com/articles/certificates/",
				"http://docs.docker.com/reference/commandline/cli/#insecure-registries",
			},
		},
		FileOwnerCheck: &FileOwnerCheck{
			filepath: "/etc/docker/certs.d/",
		},
	}
}
