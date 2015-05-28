package batten

func (dc *DockerRemoveNonEssentialSvcsCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerRemoveNonEssentialSvcsCheck) AuditCheck() (bool, error) {
	// TODO: implement
	return true, nil
}

type DockerRemoveNonEssentialSvcsCheck struct {
	*CheckDefinitionImpl
}

func makeDockerRemoveNonEssentialSvcsCheck() Check {
	return &DockerRemoveNonEssentialSvcsCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:   "CIS-Docker-Benchmark-1.5",
			category:     "Host Configuration",
			name:         "Remove all non-essential services from the host",
			impact:       "None",
			description:  "Ensure that the host running the docker daemon is running only the essential services.",
			rationale:    "It is a good practice to implement only one primary function per server to prevent functions that require different security levels from co-existing on the same server. Additionally, mixing various application environments on the same machine may hinder the granular administration of the respective applications.",
			defaultValue: `N/A`,
			auditDescription: `Inspect the docker host and ensure that it is exclusively used for running docker containers. Other services must not be found. Examples of other services include Web Server, database, or any function other than the container's main process. Some of the things you can do to inspect the system for other services are as below:

Check running processes:
ps -ef

Check socket information:
ss -nlp
or
netstat -nlp

Check inventory:
rpm -qa (or equivalent)`,
			references:  []string{"http://blog.docker.com/2013/08/containers-docker-how-secure-are-they/"},
			remediation: "Move all other services within containers controlled by Docker or to other systems.",
		},
	}
}
