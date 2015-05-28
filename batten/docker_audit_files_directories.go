package batten

import "strings"

func (dc *DockerAuditFilesDirectoriesCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

// TODO: there should be 2 types of checks: auditctl check,
// and if that fails, use a audit config file.
func (dc *DockerAuditFilesDirectoriesCheck) AuditCheck() (bool, error) {
	str, err := runAuditCtl()
	if err != nil {
		return false, err
	}

	if !PathExists(dc.path) {
		// TODO: maybe have a way to return 'not applicable?' to check results?
		return true, nil
	}

	for _, line := range strings.Split(str, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, dc.ruleCheck) {
			return true, nil
		}
	}

	return false, nil
}

type DockerAuditFilesDirectoriesCheck struct {
	*CheckDefinitionImpl
	path      string
	ruleCheck string
}

func newDockerAuditFilesDirectoriesCheckForPath(path string, id string, isFile bool) Check {

	fileStr := "file"
	if !isFile {
		fileStr = "directory"
	}

	return &DockerAuditFilesDirectoriesCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:   "CIS-Docker-Benchmark-" + id,
			category:     "Host Configuration",
			name:         "Audit Docker files and directories - " + path,
			description:  "Audit " + path,
			impact:       "Auditing generates quite big log files. Ensure to rotate and archive them periodically. Also, create a separate partition of audit to avoid filling root file system",
			rationale:    "Apart from auditing your regular Linux file system and system calls, audit all Docker related files and directories. Docker daemon runs with 'root' privileges. Its behavior depends on some key files and directories. " + path + " is one such " + fileStr + ". It must be audited.",
			defaultValue: "By default, Docker related files and directories are not audited.",
			auditDescription: `Verify that there is an audit rule corresponding to ` + path + ` ` + fileStr + `. For example, execute below command:

auditctl -l | grep ` + path + ` 

This should list a rule for ` + path + ` ` + fileStr + `.`,
			references: []string{"https://access.redhat.com/documentation/en- US/Red_Hat_Enterprise_Linux/6/html/Security_Guide/chap-system_auditing.html"},
			remediation: `Add a rule for ` + path + ` ` + fileStr + `. For example,
Add the line as below in /etc/audit/audit.rules file: 

-w ` + path + ` -k docker

Then, restart the audit daemon. For example,

#> service auditd restart`,
		},
		path:      path,
		ruleCheck: "-w " + path + " -k docker",
	}
}

func makeDockerAuditFilesVarLibDocker() Check {
	return newDockerAuditFilesDirectoriesCheckForPath("/var/lib/docker", "1.9", false)
}

func makeDockerAuditFilesEtcDocker() Check {
	return newDockerAuditFilesDirectoriesCheckForPath("/etc/docker", "1.10", false)
}

func makeDockerAuditFilesDockerRegistry() Check {
	result := newDockerAuditFilesDirectoriesCheckForPath("/usr/lib/systemd/system/docker-registry.service", "1.11", true)
	return result
}

func makeDockerAuditFilesDockerService() Check {
	result := newDockerAuditFilesDirectoriesCheckForPath("/usr/lib/systemd/system/docker.service", "1.12", true)
	return result
}

func makeDockerAuditFilesDockerSock() Check {
	result := newDockerAuditFilesDirectoriesCheckForPath("/var/run/docker.sock", "1.13", true)
	return result
}

func makeDockerAuditFilesSysconfigDocker() Check {
	result := newDockerAuditFilesDirectoriesCheckForPath("/etc/sysconfig/docker", "1.14", true)
	return result
}

func makeDockerAuditFilesSysconfigDockerNetwork() Check {
	result := newDockerAuditFilesDirectoriesCheckForPath("/etc/sysconfig/docker-network", "1.15", true)
	return result
}

func makeDockerAuditFilesSysconfigDockerRegistry() Check {
	result := newDockerAuditFilesDirectoriesCheckForPath("/etc/sysconfig/docker-registry", "1.16", true)
	return result
}
func makeDockerAuditFilesSysconfigDockerStorage() Check {
	result := newDockerAuditFilesDirectoriesCheckForPath("/etc/sysconfig/docker-storage", "1.17", true)
	return result
}

func makeDockerAuditFilesEtcDefaultDocker() Check {
	result := newDockerAuditFilesDirectoriesCheckForPath("/etc/default/docker", "1.18", true)
	return result
}
