package batten

import "os"

func (dc *DockerEnvFilePermsCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerEnvFilePermsCheck) AuditCheck() (bool, error) {
	if dc.targetPerms == 0 {
		dc.targetPerms = 0644
	}

	if PathExists(dc.filepath) {
		// TODO log actual perms for debugging
		return dc.HasAtLeastPerms(os.FileMode(dc.targetPerms))
	}

	return true, nil
}

type DockerEnvFilePermsCheck struct {
	*CheckDefinitionImpl
	*FilePermsCheck
}

func makeDockerEnvFilePermsCheck() Check {
	return &DockerEnvFilePermsCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-3.8",
			category:    `Docker daemon configuration files`,
			name:        `Verify that Docker environment file permissions are set to 644 or more restrictive`,
			description: `Docker daemon leverages Docker environment file for setting Docker daemon run time environment. If you are using Docker on a machine that uses systemd to manage services, then the file is /etc/sysconfig/docker. On other systems, the environment file is /etc/default/docker. Verify that the environment file permissions are correctly set to '644' or more restrictive.`,
			rationale:   `Docker environment file contains sensitive parameters that may alter the behavior of Docker daemon during run time. Hence, it should be only writable by 'root' to maintain the integrity of the file.`,
			auditDescription: `Execute the below command to verify that the environment file permissions are set to '644' or more restrictive:

stat -c %a <Environment file name>

For example,

stat -c %a /etc/sysconfig/docker`,
			remediation: `#> chmod 644 <Environment file name>
			￼
For example,

#> chmod 644 /etc/sysconfig/docker

This would set the file permissions for the environment file to '644'.`,
			impact:       `None.`,
			defaultValue: `By default, the file permissions for this file is correctly set to '644'.`,
			references: []string{
				"https://docs.docker.com/articles/systemd/",
			},
		},
		FilePermsCheck: &FilePermsCheck{
			filepath: "/etc/sysconfig/docker",
		},
	}

}
