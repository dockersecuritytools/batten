package batten

func (dc *DockerSystemdSocketOwnerCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerSystemdSocketOwnerCheck) AuditCheck() (bool, error) {

	if PathExists(dc.filepath) {
		return dc.IsOwnerAndGroupOwner(0, 0)
	}

	return true, nil
}

type DockerSystemdSocketOwnerCheck struct {
	*CheckDefinitionImpl
	*FileOwnerCheck
}

func makeDockerSystemdSocketOwnerCheck() Check {
	return &DockerSystemdSocketOwnerCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-3.5",
			category:    `Docker daemon configuration files`,
			name:        `Verify that docker.socket file ownership is set to root:root`,
			description: `If you are using Docker on a machine that uses systemd to manage services, then verify that the 'docker.socket' file ownership and group-ownership is correctly set to 'root'.`,
			rationale:   `'docker.socket' file contains sensitive parameters that may alter the behavior of Docker remote API. Hence, it should be owned and group-owned by 'root' to maintain the integrity of the file.`,
			auditDescription: `Execute the below command to verify that the file is owned and group-owned by 'root': 
	
stat -c %U:%G /usr/lib/systemd/system/docker.socket | grep -v root:root

The above command should not return anything.`,
			remediation: `#> chown root:root /usr/lib/systemd/system/docker.socket

This would set the ownership and group-ownership for the file to 'root'.`,
			impact:       `None.`,
			defaultValue: `This file may not be present on the system. In that case, this recommendation is not applicable. By default, if the file is present, the ownership and group-ownership for this file is correctly set to 'root'.`,
			references: []string{
				"https://docs.docker.com/articles/basics/#bind-docker-to-another-hostport-or-a- unix-socket",
				"https://github.com/YungSang/fedora-atomic-packer/blob/master/oem/docker.socket",
				"http://daviddaeschler.com/2014/12/14/centos-7rhel-7-and-docker-containers- on-boot/",
			},
		},
		FileOwnerCheck: &FileOwnerCheck{
			filepath: "/usr/lib/systemd/system/docker.socket",
		},
	}
}
