package batten

import "os"

func (dc *DockerRegistryEnvFilePermsCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerRegistryEnvFilePermsCheck) AuditCheck() (bool, error) {
	if dc.targetPerms == 0 {
		dc.targetPerms = 0644
	}

	if PathExists(dc.filepath) {
		// TODO log actual perms for debugging
		return dc.HasAtLeastPerms(os.FileMode(dc.targetPerms))
	}

	return true, nil
}

type DockerRegistryEnvFilePermsCheck struct {
	*CheckDefinitionImpl
	*FilePermsCheck
}

func makeDockerRegistryEnvFilePermsCheck() Check {
	return &DockerRegistryEnvFilePermsCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-3.12",
			category:    `Docker daemon configuration files`,
			name:        `Verify that docker-registry environment file permissions are set to 644 or more restrictive`,
			description: `If you are using Docker on a machine that uses systemd to manage services, then verify that the 'docker-registry' file permissions are correctly set to '644' or more restrictive.`,
			rationale:   `'docker-registry' file contains sensitive parameters that may alter the behavior of Docker daemon. Hence, it should not be writable by any other user other than 'root' to maintain the integrity of the file.`,
			auditDescription: `Execute the below command to verify that the file permissions are set to '644' or more restrictive:

stat -c %a /etc/sysconfig/docker-registry`,
			remediation: `#> chmod 644 /etc/sysconfig/docker-registry

This would set the file permissions to '644'.`,
			impact:       `None`,
			defaultValue: `This file may not be present on the system. In that case, this recommendation is not applicable. By default, if the file is present, the file permissions are correctly set to '644'.`,
			references: []string{
				"https://docs.docker.com/articles/systemd/",
				"",
			},
		}, FilePermsCheck: &FilePermsCheck{
			filepath: "/etc/sysconfig/docker-registry",
		},
	}
}
