package batten

func (dc *DockerCheckEndpointProtectionPlatform) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerCheckEndpointProtectionPlatform) AuditCheck() (bool, error) {
	// TODO
	return true, nil
}

type DockerCheckEndpointProtectionPlatform struct {
	*CheckDefinitionImpl
}

func makeDockerCheckEndpointProtectionPlatform() Check {
	return &DockerCheckEndpointProtectionPlatform{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-6.3",
			category:    "Docker Security Operations",
			name:        `Endpoint protection platform (EPP) tools for containers`,
			description: `There is no container-aware endpoint protection platform (EPP) solution as of now. You must rely on compensating controls to achieve the same.`,
			rationale:   `Traditional EPP and encryption vendors have not yet recognized containers as an area that they need to pursue and secure in the future. Hence, there are no suitable products at this time. Thus, you must rely on compensating controls.`,
			auditDescription: `Verify whether compensating controls, in place of traditional EPP, exist. Some of the compensating controls might be as below: 
	•	Application white-listing (as supported by AppArmor)  
	•	Mandatory access controls (as supported by SELinux) or  
	•	Make containers self-assessing entities using the DevOps tool chain`,
			remediation:  `AppArmor, SELinux and DevOps product configurations for containers are beyond the scope of this benchmark. You should seek guidance on specific configuration needed for containers from their respective sources.`,
			impact:       "None",
			defaultValue: `By default, no endpoint protection is provided to containers.`,
			references: []string{
				"https://www.gartner.com/doc/2956826/security-properties-containers-managed-docker",
			},
		},
	}
}
