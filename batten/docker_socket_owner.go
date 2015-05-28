package batten

import "errors"

func (dc *DockerSocketOwnerCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerSocketOwnerCheck) AuditCheck() (bool, error) {
	if PathExists(dc.filepath) {
		return dc.validateOwnerAndGroupOwner()
	} else {
		return false, errors.New("Could not find path: " + dc.filepath)
	}
	return true, nil
}

type DockerSocketOwnerCheck struct {
	*CheckDefinitionImpl
	*FileOwnerCheck
}

func makeDockerSocketOwnerCheck() Check {
	return &DockerSocketOwnerCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-3.25",
			category:    `Docker daemon configuration files`,
			name:        `Verify that Docker socket file ownership is set to root:docker`,
			description: `Verify that the Docker socket file is owned by 'root' and group-owned by 'docker'.`,
			rationale: `Docker daemon runs as 'root'. The default Unix socket hence must be owned by 'root'. If any other user or process owns this socket, then it might be possible for that non- privileged user or process to interact with Docker daemon. Also, such a non-privileged user or process might interact with containers. This is neither secure nor desired behavior.

Additionally, the Docker installer creates a Unix group called 'docker'. You can add users to this group, and then those users would be able to read and write to default Docker Unix socket. The membership to the 'docker' group is tightly controlled by the system administrator. If any other group owns this socket, then it might be possible for members of that group to interact with Docker daemon. Also, such a group might not be as tightly controlled as the 'docker' group. This is neither secure nor desired behavior.

Hence, the default Docker Unix socket file must be owned by 'root' and group-owned by 'docker' to maintain the integrity of the socket file.`,
			auditDescription: `Execute the below command to verify that the Docker socket file is owned by 'root' and group-owned by 'docker':

stat -c %U:%G /var/run/docker.sock | grep -v root:docker
The above command should not return anything.`,
			remediation: `#> chown root:docker /var/run/docker.sock

This would set the ownership to 'root' and group-ownership to 'docker' for default Docker socket file.`,
			impact:       `None`,
			defaultValue: `By default, the ownership and group-ownership for Docker socket file is correctly set to 'root:docker'.`,
			references: []string{
				"https://docs.docker.com/reference/commandline/cli/#daemon-socket-option",
				"https://docs.docker.com/articles/basics/#bind-docker-to-another-hostport-or-a-unix-socket",
			},
		},
		FileOwnerCheck: &FileOwnerCheck{username: "root", groupname: "docker", filepath: "/var/run/docker.sock"},
	}
}
