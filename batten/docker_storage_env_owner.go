package batten

func (dc *DockerStorageEnvOwnerCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerStorageEnvOwnerCheck) AuditCheck() (bool, error) {
	if PathExists(dc.filepath) {
		return dc.IsOwnerAndGroupOwner(0, 0)
	}

	return true, nil
}

type DockerStorageEnvOwnerCheck struct {
	*CheckDefinitionImpl
	*FileOwnerCheck
}

func makeDockerStorageEnvOwnerCheck() Check {
	return &DockerStorageEnvOwnerCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-3.13",
			category:    `Docker daemon configuration files`,
			name:        `Verify that docker-storage environment file ownership is set to root:root`,
			description: `If you are using Docker on a machine that uses systemd to manage services, then verify that the 'docker-storage' file ownership and group-ownership is correctly set to 'root'.`,
			rationale:   `'docker-storage' file contains sensitive parameters that may alter the behavior of Docker daemon. Hence, it should be owned and group-owned by 'root' to maintain the integrity of the file.`,
			auditDescription: `Execute the below command to verify that the file is owned and group-owned by 'root': 
			
stat -c %U:%G /etc/sysconfig/docker-storage | grep -v root:root

The above command should not return anything.`,
			remediation: `#> chown root:root /etc/sysconfig/docker-storage
			￼￼
This would set the ownership and group-ownership for the file to 'root'.`,
			impact:       `None`,
			defaultValue: `This file may not be present on the system. In that case, this recommendation is not applicable. By default, if the file is present, the ownership and group-ownership for this file is correctly set to 'root'.`,
			references: []string{
				"https://docs.docker.com/articles/systemd/",
			},
		},
		FileOwnerCheck: &FileOwnerCheck{
			filepath: "/etc/sysconfig/docker-storage",
		},
	}
}
