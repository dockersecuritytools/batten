package batten

func (dc *DockerInsecureRegistriesCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerInsecureRegistriesCheck) AuditCheck() (bool, error) {
	succ, environ, err := readDockerDaemonEnviron(dc.dockerPidFile)
	if err != nil {
		return false, err
	}
	if succ && environ["insecure-registry"] != "" {
		return false, nil
	}
	return true, nil
}

type DockerInsecureRegistriesCheck struct {
	*CheckDefinitionImpl
	// TODO: make configurable
	dockerPidFile string
}

func makeDockerInsecureRegistriesCheck() Check {
	return &DockerInsecureRegistriesCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-2.5",
			category:    `Docker daemon configuration`,
			name:        `Do not use insecure registries`,
			description: `Docker considers a private registry either secure or insecure. By default, registries are considered secure.`,
			rationale: `A secure registry uses TLS. A copy of registry's CA certificate is placed on the Docker host at '/etc/docker/certs.d/<registry-name>/' directory. An insecure registry is the one not having either valid registry certificate or is not using TLS. You should not be using any insecure registries in the production environment. Insecure registries can be tampered with leading to possible compromise to your production system.
			Additionally, If a registry is marked as insecure then 'docker pull', 'docker push', and 'docker search' commands will not result in an error message and the user might be indefinitely working with insecure registries without ever being notified of potential danger.`,
			auditDescription: `xecute the below command to find out if any insecure registries are used:

$> ps -ef | grep docker
			￼￼￼
Ensure that the '--insecure-registry' parameter is not present.`,
			remediation: `Do not use any insecure registries.

For example, do not start the Docker daemon as below: 

$> docker -d --insecure-registry 10.1.0.0/16`,
			impact:       `None`,
			defaultValue: `By default, Docker assumes all, but local, registries are secure.`,
			references: []string{
				"http://docs.docker.com/reference/commandline/cli/#insecure-registries",
			},
		},
		dockerPidFile: "/var/run/docker.pid",
	}
}
