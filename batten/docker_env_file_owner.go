package batten

func (dc *DockerEnvFileOwnerCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerEnvFileOwnerCheck) AuditCheck() (bool, error) {
	if PathExists(dc.filepath) {
		return dc.IsOwnerAndGroupOwner(0, 0)
	}

	return true, nil
}

type DockerEnvFileOwnerCheck struct {
	*CheckDefinitionImpl
	*FileOwnerCheck
}

func makeDockerEnvFileOwnerCheck() Check {
	return &DockerEnvFileOwnerCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-3.7",
			category:    `Docker daemon configuration files`,
			name:        `Verify that Docker environment file ownership is set to root:root`,
			description: `Docker daemon leverages Docker environment file for setting Docker daemon run time environment. If you are using Docker on a machine that uses systemd to manage services, then the file is /etc/sysconfig/docker. On other systems, the environment file is /etc/default/docker. Verify that the environment file ownership and group-ownership is correctly set to 'root'.`,
			rationale:   `Docker environment file contains sensitive parameters that may alter the behavior of Docker daemon during run time. Hence, it should be owned and group-owned by 'root' to maintain the integrity of the file.`,
			auditDescription: `Execute the below command to verify that the environment file is owned and group-owned by 'root':

stat -c %U:%G <Environment file name> | grep -v root:root
   
For example,

stat -c %U:%G /etc/sysconfig/docker | grep -v root:root

The above command should not return anything.`,
			remediation: `#> chown root:root <Environment file name>
For example,

#> chown root:root /etc/sysconfig/docker

This would set the ownership and group-ownership for the environment file to 'root'.`,
			impact:       `None`,
			defaultValue: `By default, the ownership and group-ownership for this file is correctly set to 'root'.`,
			references: []string{
				"https://docs.docker.com/articles/systemd/",
			},
		},

		FileOwnerCheck: &FileOwnerCheck{
			filepath: "/etc/sysconfig/docker",
		},
	}
}
