package batten

func (dc *DockerXXX) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerXXX) AuditCheck() (bool, error) {
	// TODO: implement
	return true, nil
}

type DockerXXX struct {
	*CheckDefinitionImpl
}

func makeDockerXXX() Check {
	return &DockerXXX{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:       "CIS-Docker-Benchmark-???",
			category:         ``,
			name:             ``,
			description:      ``,
			rationale:        ``,
			auditDescription: ``,
			remediation:      ``,
			impact:           ``,
			defaultValue:     ``,
			references: []string{
				"",
				"",
				"",
			},
		},
	}
}
