package batten

func (dc *DockerMonitorContainers) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerMonitorContainers) AuditCheck() (bool, error) {
	// TODO
	return true, nil
}

type DockerMonitorContainers struct {
	*CheckDefinitionImpl
}

func makeDockerMonitorContainers() Check {
	return &DockerMonitorContainers{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-6.2",
			category:    "Docker Security Operations",
			name:        `Monitor Docker containers usage, performance and metering`,
			description: `Containers might run services that are critical for your business. Monitoring their usage, performance and metering would be of paramount importance.`,
			rationale: `Tracking container usage, performance and having some sort of metering around them would be important as you embrace the containers to run critical services for your business. This would give you 
	•	Capacity Management and Optimization  
	•	Performance Management  
	•	Comprehensive Visibility  Such a deep visibility of container performance would help you ensure high availability of containers and minimum downtime.`,
			auditDescription: `Verify whether a container or a software is used for tracking container usage, reporting performance and metering.`,
			remediation:      `Use a software or a container for tracking container usage, reporting performance and metering.`,
			impact:           "None",
			defaultValue:     `By default, for each container, runtime metrics about CPU, memory, and block I/O usage is tracked by the system via enforcement of control groups (cgroups) as below:  CPU - /sys/fs/cgroup/cpu/system.slice/docker-$INSTANCE_ID.scope/ Memory - /sys/fs/cgroup/memory/system.slice/docker-$INSTANCE_ID.scope/  Block I/O - /sys/fs/cgroup/blkio/system.slice/docker-$INSTANCE_ID.scope/`,
			references: []string{
				"https://docs.docker.com/articles/runmetrics/",
				"https://github.com/google/cadvisor",
				"https://docs.docker.com/reference/commandline/cli/#stats",
			},
		},
	}
}
