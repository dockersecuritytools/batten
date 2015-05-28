package batten

import (
	"strings"

	"github.com/fsouza/go-dockerclient"
)

func (dc *DockerSingleMainProcess) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerSingleMainProcess) AuditCheck() (bool, error) {
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

		data, err := client.TopContainer(c.ID, "-el")

		if err != nil {
			continue
		}

		if len(data.Processes) > 1 {
			return false, nil
		}

		// Makre sure the main process is not a
		// process monitoring application
		last := data.Processes[0][len(data.Processes[0])-1]
		for _, item := range dc.processManagers {
			if strings.Contains(last, item) {
				return false, nil
			}
		}

	}

	return true, nil
}

type DockerSingleMainProcess struct {
	*CheckDefinitionImpl
	processManagers []string
}

func makeDockerSingleMainProcess() Check {
	return &DockerSingleMainProcess{
		processManagers: []string{
			"supervisord",
			"monitd",
		},
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-5.3",
			category:    `Container Runtime`,
			name:        `Verify that containers are running only a single main process`,
			description: `In almost all cases, you should only run a single main process (that main process could spawn children, which is ok) in a single container. Decoupling applications into multiple containers makes it much easier to scale horizontally and reuse containers. If that service depends on another service, make use of container linking.`,
			rationale: `By design, Docker watches one single process within the container. So, installing and running multiple applications within a single container breaks the basic design of 'one container one process'. 
If you need multiple processes, you need to add one at the top-level to take care of the others. You also need to add a process manager; for instance Monit or Supervisor. In other words, you're turning a lean and simple container into something much more complicated. If your application stops (if it exits cleanly or if it crashes), instead of getting that information through Docker, you will have to get it from your process manager.`,
			auditDescription: `Step 1: List all the running instances of containers by executing below command: 
docker ps -q 
Step 2: For each container instance, execute the below command: 
docker exec $INSTANCE_ID ps -el 
Ensure that there is only one process running i.e. the process the container is intended to run. There should not be a process for process manager such as Monit, Supervisor or any other.`,
			remediation:  `Do not run multiple applications within a single container. Use container linking instead to run multiple applications in multiple containers in tandem.`,
			impact:       `None`,
			defaultValue: `By default, only one process per container is allowed.`,
			references: []string{
				"https://docs.docker.com/articles/dockerfile_best-practices",
				"https://docs.docker.com/userguide/dockerlinks",
				"https://docs.docker.com/articles/using_supervisord",
			},
		},
	}
}
