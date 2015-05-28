package batten

import docker "github.com/fsouza/go-dockerclient"

func (dc *DockerContainerUserCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

// list all running containers, and ensure they are all running as root
func (dc *DockerContainerUserCheck) AuditCheck() (bool, error) {

	client, err := docker.NewClient(DockerUnixSocket)

	if err != nil {
		// TODO: log error message
		return false, err
	}

	containers, err := client.ListContainers(docker.ListContainersOptions{All: false})

	if err != nil {
		// TODO: log error message
		return false, err
	}
	for _, c := range containers {
		container, err := client.InspectContainer(c.ID)
		if err != nil {
			// TODO: log error message
			return false, err
		}

		if container.Config != nil && container.Config.User == "" {
			// TODO log the container with the bad user?
			return false, nil
		}

	}
	return true, nil
}

type DockerContainerUserCheck struct {
	*CheckDefinitionImpl
	// TODO: allow the user to specify a policy of containers they are ok with running as root.
}

func makeDockerContainerUserCheck() Check {
	return &DockerContainerUserCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-4.1",
			category:    `Container Images and Build File`,
			name:        `Create a user for the container`,
			description: `Create a non-root user for the container in the Dockerfile for the container image. Also, run the container with non-root user.`,
			rationale:   `Currently, mapping the container's root user to a non-root user on the host is not supported by Docker. The support for user namespace would be provided in future releases (probably in 1.6). This creates a serious user isolation issue. It is thus highly recommended to ensure that there is a non-root user created for the container and the container is run using that user.`,
			auditDescription: `docker ps -q | xargs docker inspect --format '{{ .Id }}: User={{.Config.User}}'

The above command should return container username or user ID. If it is blank it means, the container is running as root.'`,
			remediation: `Ensure that the Dockerfile for the container image contains below instruction:

USER <username or ID>

where username or ID refers to the user that could be found in the container base image. If there is no specific user created in the container base image, then add a useradd command to add the specific user before USER instruction.

For example, add the below lines in the Dockerfile to create a user in the container:

When you run the container, use the '-u' flag to specify that you would want to run the container as a specific user and not root. This can be done by executing below command:

$> docker run -u <Username or ID> <Run args> <Container Image Name or ID> <Command>

For example,

$> docker run -u 1000 -i -t centos /bin/bash

This would ensure that the above container is run with user ID 1000 and not root.
Note: If there are users in the image that the containers do not need, consider deleting them. After deleting those users, commit the image and then generate new instances of containers for use.`,
			impact:       `None.`,
			defaultValue: `By default, the containers are run with root privileges and as user root inside the container.`,
			references: []string{
				"https://github.com/docker/docker/issues/2918",
				"https://github.com/docker/docker/pull/4572",
				"https://github.com/docker/docker/issues/7906",
				"https://www.altiscale.com/hadoop-blog/making-docker-work-yarn/",
				"http://docs.docker.com/articles/security/",
			},
		},
	}
}
