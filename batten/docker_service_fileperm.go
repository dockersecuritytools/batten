package batten

import "os"

func (dc *DockerSvcFilePermsCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerSvcFilePermsCheck) AuditCheck() (bool, error) {

	if dc.targetPerms == 0 {
		dc.targetPerms = 0644
	}

	if PathExists(dc.filepath) {
		return dc.HasAtLeastPerms(os.FileMode(dc.targetPerms))
	}

	return true, nil
}

type DockerSvcFilePermsCheck struct {
	*CheckDefinitionImpl
	*FilePermsCheck
}

func makeDockerSvcFilePermsCheck() Check {
	return &DockerSvcFilePermsCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-3.2",
			category:    `Docker daemon configuration files`,
			name:        `Verify that docker.service file permissions are set to 644 or more restrictive`,
			description: `If you are using Docker on a machine that uses systemd to manage services, then verify that the 'docker.service' file permissions are correctly set to '644' or more restrictive.`,
			rationale:   `'docker.service' file contains sensitive parameters that may alter the behavior of Docker daemon. Hence, it should not be writable by any other user other than 'root' to maintain the integrity of the file.`,
			auditDescription: `Execute the below command to verify that the file permissions are set to '644' or more restrictive:

stat -c %a /usr/lib/systemd/system/docker.service`,
			remediation: `#> chmod 644 /usr/lib/systemd/system/docker.service
This would set the file permissions to '644'.`,
			impact:       `None.`,
			defaultValue: `This file may not be present on the system. In that case, this recommendation is not applicable. By default, if the file is present, the file permissions are correctly set to '644'.`,
			references: []string{
				"https://docs.docker.com/articles/systemd",
			},
		},
		FilePermsCheck: &FilePermsCheck{
			filepath: "/usr/lib/systemd/system/docker.service",
		},
	}
}
