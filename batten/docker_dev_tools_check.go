package batten

func (dc *DockerDevToolsCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerDevToolsCheck) AuditCheck() (bool, error) {
	// TODO: define a policy of tools that are not allowed
	// to be installed or present in memory on the host container
	return true, nil
}

type DockerDevToolsCheck struct {
	*CheckDefinitionImpl
	devToolsPolicy []string
}

func makeDockerDevToolsCheck() Check {
	return &DockerDevToolsCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:       "CIS-Docker-Benchmark-1.3",
			name:             "Do not use development tools in production",
			category:         "Host Configuration",
			impact:           "None",
			description:      "Development tools are not meant to be used in production.",
			rationale:        "Development tools usually provide alternate and convenient methods of using a technology. These tools are just meant to develop quick proof-of-concept or provide leads for up-selling of production ready software. There are various development tools that can be used with Docker such as boot2docker, Kitematic, VMware Fusion, Vagrant and others. Do not use these tools in production environment.",
			defaultValue:     `N/A`,
			auditDescription: "Inspect the environment that is running Docker and verify that development tools are not used in production.",
			references:       []string{"https://github.com/boot2docker/boot2docker", "https://kitematic.com/"},
			remediation:      `Do not use development tools in production.`,
		},
	}
}
