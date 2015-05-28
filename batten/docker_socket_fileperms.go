package batten

import "os"

func (dc *DockerSocketFilePermsCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerSocketFilePermsCheck) AuditCheck() (bool, error) {

	if PathExists(dc.filepath) {
		// TODO log actual perms for debugging
		return dc.HasAtLeastPerms(os.FileMode(dc.targetPerms))
	}

	return true, nil
}

type DockerSocketFilePermsCheck struct {
	*CheckDefinitionImpl
	*FilePermsCheck
}

func makeDockerSocketFilePermsCheck() Check {
	return &DockerSocketFilePermsCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-3.26",
			category:    `Docker daemon configuration files`,
			name:        `Verify that Docker socket file permissions are set to 660 or more restrictive`,
			description: `Verify that the Docker socket file has permissions of '660' or more restrictive.`,
			rationale:   `Only 'root' and members of 'docker' group should be allowed to read and write to default Docker Unix socket. Hence, the Docket socket file must have permissions of '660' or more restrictive.`,
			auditDescription: `Execute the below command to verify that the Docker socket file has permissions of '660' or more restrictive:

stat -c %a /var/run/docker.sock`,
			remediation: `#> chmod 660 /var/run/docker.sock

This would set the file permissions of the Docker socket file to '660'.`,
			impact:       `None`,
			defaultValue: `By default, the permissions for Docker socket file is correctly set to '660'.`,
			references: []string{
				"https://docs.docker.com/reference/commandline/cli/#daemon-socket-option",
				"https://docs.docker.com/articles/basics/#bind-docker-to-another-hostport-or-a-unix-socket",
			},
		},
		FilePermsCheck: &FilePermsCheck{
			targetPerms: 0660,
			filepath:    "/var/run/docker.sock",
		},
	}
}
