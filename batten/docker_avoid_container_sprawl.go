package batten

import docker "github.com/fsouza/go-dockerclient"

func (dc *DockerAvoidContainerSprawl) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerAvoidContainerSprawl) AuditCheck() (bool, error) {
	client, err := getDockerAPIConnection()

	if err != nil {
		// TODO: log error message
		return false, err
	}

	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})

	if err != nil {
		// TODO: log error message
		return false, err
	}

	var count int
	for _, c := range containers {

		if cc, err := client.InspectContainer(c.ID); err == nil {
			if cc.State.Running {
				count++
			}
		}

	}

	if count < len(containers) {
		return false, nil
	}

	return true, nil
}

type DockerAvoidContainerSprawl struct {
	*CheckDefinitionImpl
}

func makeDockerAvoidContainerSprawl() Check {
	return &DockerAvoidContainerSprawl{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-6.7",
			category:    "Docker Security Operations",
			name:        "Avoid container sprawl",
			impact:      "If you keep way too few number of containers per host, then perhaps you are not utilizing your host resources very adequately.",
			description: `Do not keep a large number of containers on the same host.`,
			rationale:   `The flexibility of containers makes it easy to run multiple instances of applications and indirectly leads to Docker images that exist at varying security patch levels. It also means that you are consuming host resources that otherwise could have been used for running 'useful' containers. Having more than just the manageable number of containers on a particular host makes the situation vulnerable to mishandling, misconfiguration and fragmentation. Thus, avoid container sprawl and keep the number of containers on a host to a manageable total.`,
			auditDescription: `Execute the below command to find the total number of containers you have on the host: 
docker info | grep "Containers" 
Now, execute the below command to find the total number of containers that are actually running on the host. 
docker ps -q | wc -l 
If the difference between the number of containers that are present on the host and the number of containers that are actually running on the host is large (say 25 or more), then perhaps, the containers are sprawled on the host.`,
			remediation: `Periodically check your container inventory per host and clean up the containers that are not needed using the below command: 
$> docker rm $INSTANCE_ID 
For example, 
$> docker rm e3a7a1a97c58`,
			defaultValue: `By default, Docker does not restrict the number of containers you may have on a host.`,
			references: []string{
				"https://zeltser.com/security-risks-and-benefits-of-docker-application/",
				"http://searchsdn.techtarget.com/feature/Docker-networking-How-Linux-containers-will-change-your-network",
			},
		},
	}
}
