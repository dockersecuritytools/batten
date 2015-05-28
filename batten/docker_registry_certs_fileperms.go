package batten

import "os"

func (dc *DockerRegistryCertsFilePermsCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerRegistryCertsFilePermsCheck) AuditCheck() (bool, error) {
	if dc.targetPerms == 0 {
		dc.targetPerms = 0444
	}

	if PathExists(dc.filepath) {
		// TODO log actual perms for debugging
		return dc.HasAtLeastPermsRecursive(os.FileMode(dc.targetPerms))
	}

	return true, nil
}

type DockerRegistryCertsFilePermsCheck struct {
	*CheckDefinitionImpl
	*FilePermsCheck
}

func makeDockerRegistryCertsFilePermsCheck() Check {
	return &DockerRegistryCertsFilePermsCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-3.18",
			category:    `Docker daemon configuration files`,
			name:        `Verify that registry certificate file permissions are set to 444 or more restrictive`,
			description: `Verify that all the registry certificate files (usually found under /etc/docker/certs.d/<registry-name> directory) have permissions of '444' or more restrictive.`,
			rationale:   `/etc/docker/certs.d/<registry-name> directory contains Docker registry certificates. These certificate files must have permissions of '444' to maintain the integrity of the certificates.`,
			auditDescription: `Execute the below command to verify that the registry certificate files have permissions of '444' or more restrictive:

stat -c %a /etc/docker/certs.d/<registry-name>/*`,
			remediation: `#> chmod 444 /etc/docker/certs.d/<registry-name>/*

This would set the permissions for registry certificate files to '444'.`,
			impact:       `None`,
			defaultValue: `By default, the permissions for registry certificate files might not be '444'. The default file permissions are governed by the system or user specific umask values.`,
			references: []string{
				"https://docs.docker.com/articles/certificates/",
				"http://docs.docker.com/reference/commandline/cli/#insecure-registries",
			},
		},
		FilePermsCheck: &FilePermsCheck{
			filepath: "/etc/docker/certs.d/",
		},
	}
}
