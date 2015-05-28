package batten

import "strings"

func (dc *DockerDaemonAuditingCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerDaemonAuditingCheck) checkUsingAuditctl() (bool, error) {

	str, err := runAuditCtl()
	if err != nil {
		return false, err
	}

	for _, line := range strings.Split(str, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, dc.ruleCheck) {
			return true, nil
		}
	}

	return false, nil
}

func (dc *DockerDaemonAuditingCheck) AuditCheck() (bool, error) {

	return dc.checkUsingAuditctl()
}

type DockerDaemonAuditingCheck struct {
	*CheckDefinitionImpl
	auditRulesFile string
	// TODO: make this configureable
	auditCtlPath string
	ruleCheck    string
}

func makeDockerDaemonAuditingCheck() Check {
	return &DockerDaemonAuditingCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:   "CIS-Docker-Benchmark-1.8",
			category:     "Host Configuration",
			name:         "Audit docker daemon",
			impact:       "Auditing generates quite big log files. Ensure to rotate and archive them periodically. Also, create a separate partition of audit to avoid filling root file system.",
			description:  "Audit all Docker daemon activities.",
			rationale:    "Apart from auditing your regular Linux file system and system calls, audit Docker daemon as well. Docker daemon runs with 'root' privileges. It is thus necessary to audit its activities and usage.",
			defaultValue: `By default, Docker daemon is not audited.`,
			auditDescription: `Verify that there is an audit rule for Docker daemon. For example, execute below command:

auditctl -l | grep /usr/bin/docker

This should list a rule for Docker daemon.`,
			references: []string{
				"https://access.redhat.com/documentation/en- US/Red_Hat_Enterprise_Linux/6/html/Security_Guide/chap-system_auditing.html",
			},
			remediation: `Add a rule for Docker daemon.
For example, Add the line as below line in /etc/audit/audit.rules file:

 -w /usr/bin/docker -k docker

Then, restart the audit daemon. For example,

#> service auditd restart`,
		},
		ruleCheck: "-w /usr/bin/docker -k docker",
	}
}
