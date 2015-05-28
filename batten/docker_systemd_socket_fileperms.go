package batten

import "os"

func (dc *DockerSystemdSocketFilePermsCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerSystemdSocketFilePermsCheck) AuditCheck() (bool, error) {
	if dc.targetPerms == 0 {
		dc.targetPerms = 0644
	}

	if PathExists(dc.filepath) {
		// TODO log actual perms for debugging
		return dc.HasAtLeastPerms(os.FileMode(dc.targetPerms))
	}

	return true, nil
}

type DockerSystemdSocketFilePermsCheck struct {
	*CheckDefinitionImpl
	*FilePermsCheck
}

func makeDockerSystemdSocketFilePermsCheck() Check {
	return &DockerSystemdSocketFilePermsCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-3.6",
			category:    `Docker daemon configuration files`,
			name:        `Verify that docker.socket file permissions are set to 644 or more restrictive`,
			description: `If you are using Docker on a machine that uses systemd to manage services, then verify that the 'docker.socket' file permissions are correctly set to '644' or more restrictive.`,
			rationale:   `'docker.socket' file contains sensitive parameters that may alter the behavior of Docker remote API. Hence, it should be writable only by 'root' to maintain the integrity of the file.`,
			auditDescription: `Execute the below command to verify that the file permissions are correctly set to '644' or more restrictive:
			
stat -c %a /usr/lib/systemd/system/docker.socket`,
			remediation: `#> chmod 644 /usr/lib/systemd/system/docker.socket
This would set the file permissions for this file to '644'.`,
			impact:       `None.`,
			defaultValue: `This file may not be present on the system. In that case, this recommendation is not applicable. By default, if the file is present, the file permissions for this file are correctly set to '644'.`,
			references: []string{
				"https://docs.docker.com/articles/basics/#bind-docker-to-another-hostport-or-a- unix-socket",
				"https://github.com/YungSang/fedora-atomic-packer/blob/master/oem/docker.socket",
				"http://daviddaeschler.com/2014/12/14/centos-7rhel-7-and-docker-containers- on-boot/",
			},
		},
		FilePermsCheck: &FilePermsCheck{
			filepath: "/usr/lib/systemd/system/docker.socket",
		},
	}
}
