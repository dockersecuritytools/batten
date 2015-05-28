package batten

import "os"

func (dc *DockerEtcDockerFilePermsCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerEtcDockerFilePermsCheck) AuditCheck() (bool, error) {
	if dc.targetPerms == 0 {
		dc.targetPerms = 0755
	}

	if PathExists(dc.filepath) {
		// TODO log actual perms for debugging
		return dc.HasAtLeastPerms(os.FileMode(dc.targetPerms))
	}

	return true, nil
}

type DockerEtcDockerFilePermsCheck struct {
	*CheckDefinitionImpl
	*FilePermsCheck
}

func makeDockerEtcDockerFilePermsCheck() Check {
	return &DockerEtcDockerFilePermsCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-3.16",
			category:    `Docker daemon configuration files`,
			name:        `Verify that /etc/docker directory permissions are set to 755 or more restrictive`,
			description: `Verify that the /etc/docker directory permissions are correctly set to '755' or more restrictive.`,
			rationale:   `'/etc/docker' directory contains certificates and keys in addition to various sensitive files. Hence, it should only be writable by 'root' to maintain the integrity of the directory.`,
			auditDescription: `Execute the below command to verify that the directory has permissions of '755' or more restrictive:

stat -c %a /etc/docker`,
			remediation: `#> chmod 755 /etc/docker
			￼
This would set the permissions for the directory to '755'.`,
			impact:       `None`,
			defaultValue: `By default, the permissions for this directory are correctly set to '755'.`,
			references: []string{
				"https://docs.docker.com/articles/certificates/",
			},
		},
		FilePermsCheck: &FilePermsCheck{
			filepath: "/etc/docker",
		},
	}
}
