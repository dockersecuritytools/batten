package batten

func (dc *DockerPerformSecurityAudits) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerPerformSecurityAudits) AuditCheck() (bool, error) {
	// TODO
	return true, nil
}

type DockerPerformSecurityAudits struct {
	*CheckDefinitionImpl
}

func makeDockerPerformSecurityAudits() Check {
	return &DockerPerformSecurityAudits{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:       "CIS-Docker-Benchmark-6.1",
			category:         "Docker Security Operations",
			name:             `Perform regular security audits of your host system and containers`,
			description:      `Perform regular security audits of your host system and containers to identify any mis- configurations or vulnerabilities that could expose your system to compromise.`,
			rationale:        `Performing regular and dedicated security audits of your host systems and containers could provide deep security insights that you might not know in your daily course of business. The identified security weaknesses should be then mitigated and this overall improves security posture of your environment.`,
			auditDescription: `Follow your organization's security audit policies and requirements.`,
			remediation:      `Follow your organization's security audit policies and requirements.`,
			impact:           "None",
			defaultValue:     `Not applicable.`,
			references: []string{
				"http://searchsecurity.techtarget.com/IT-security-auditing-Best-practices-for-conducting-audits",
			},
		},
	}
}
