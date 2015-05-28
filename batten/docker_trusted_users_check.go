package batten

import (
	"errors"
	"io/ioutil"
	"strings"
)

func (dc *DockerTrustedUsersCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerTrustedUsersCheck) AuditCheck() (bool, error) {
	bytes, err := ioutil.ReadFile(dc.groupsFile)

	if err != nil {
		return false, err
	}

	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		fields := strings.Split(line, ":")

		if fields[0] == "docker" {

			if len(fields) > 1 {
				last := fields[len(fields)-1]
				users := strings.Split(last, ",")
				for _, user := range users {
					user = strings.TrimSpace(user)
					if len(user) == 0 {
						continue
					}

					if !stringInSlice(user, dc.trustedUsers) {
						return false, errors.New(user + " is not a trusted dockergroup user")
					}
				}
			}
		}

	}

	return true, nil
}

type DockerTrustedUsersCheck struct {
	*CheckDefinitionImpl
	trustedUsers []string
	groupsFile   string
}

func makeDockerTrustedUsersCheck() Check {
	return &DockerTrustedUsersCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:   "CIS-Docker-Benchmark-1.7",
			category:     "Host Configuration",
			name:         "Only allow trusted users to control Docker daemon",
			impact:       "Rights to build and execute containers as normal user would be restricted.",
			description:  "The Docker daemon currently requires 'root' privileges. A user added to the 'docker' group gives him full 'root' access rights.",
			rationale:    "Docker allows you to share a directory between the Docker host and a guest container without limiting the access rights of the container. This means that you can start a container and map the / directory on your host to the container. The container will then be able to alter your host file system without any restrictions. In simple terms, it means that you can attain elevated privileges with just being a member of the 'docker' group and then starting a container with mapped / directory on the host.",
			defaultValue: `N/A`,
			auditDescription: `Execute the below command on the docker host and ensure that only trusted users are members of the 'docker' group.
			
cat /etc/group | grep docker`,
			references: []string{
				"https://docs.docker.com/articles/security/#docker-daemon-attack-surface",
				"https://www.andreas-jung.com/contents/on-docker-security-docker-group-considered-harmful",
			},
			remediation: "Remove any users from the 'docker' group that are not trusted. Additionally, do not create a mapping of sensitive directories on host to container volumes.",
		},
		groupsFile: "/etc/group",
	}
}
