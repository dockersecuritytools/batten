package batten

func (dc *DockerCheckCentralLogCollection) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerCheckCentralLogCollection) AuditCheck() (bool, error) {
	// TODO
	return true, nil
}

type DockerCheckCentralLogCollection struct {
	*CheckDefinitionImpl
}

func makeDockerCheckCentralLogCollection() Check {
	return &DockerCheckCentralLogCollection{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-6.5",
			category:    "Docker Security Operations",
			name:        "Use a centralized and remote log collection service",
			description: `Each container maintains its logs under /var/lib/docker/containers/$INSTANCE_ID/$INSTANCE_ID-json.log. But, maintaining logs at a centralized place is preferable.`,
			rationale: `Storing log data on a remote host or a centralized place protects log integrity from local attacks. If an attacker gains access on the local system, he could tamper with or remove log data that is stored on the local system. Also, the 'docker logs' paradigm is not yet fully developed. There are quite a few difficulties in managing the container logs namely 
	•	No logrotate for container logs  
	•	Transient behavior of docker logs  
	•	Difficulty in accessing application specific log files  
	•	All stdout and stderr are logged  Hence, a centralized and remote log collection service should be utilized to keep logs for all the containers.`,
			auditDescription: `First, verify that the centralized and remote log collection service is configured. Then verify that all the containers are logging at this centralized place.  Step 1: List all the running instances of containers by executing below command: docker ps -q  Step 2: For each container instance, execute the below command:  docker inspect --format='{{.Volumes}}' $INSTANCE_ID  Ensure that a centralized log volume is mounted on the containers.`,
			remediation:      `Configure a centralized and remote log collection service. Some of the examples to do this are in references. Once the log collection service is active, configure all the containers to send their logs to this service.`,
			impact:           "None",
			defaultValue:     `By default, each container logs separately.`,
			references: []string{
				"https://docs.docker.com/reference/commandline/cli/#logs",
				"http://jpetazzo.github.io/2014/08/24/syslog-docker/",
				"http://stackengine.com/docker-logs-aggregating-ease/",
			},
		},
	}
}
