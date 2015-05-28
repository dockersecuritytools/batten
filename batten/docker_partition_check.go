package batten

import (
	"io/ioutil"
	"strings"
)

var DEFAULT_FSTAB = "/etc/fstab"

func (dc *DockerPartitionCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerPartitionCheck) AuditCheck() (bool, error) {
	bytes, err := ioutil.ReadFile(dc.fstab)

	if err != nil {
		return false, err
	}

	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)

		if len(fields) > 1 && fields[1] == "/var/lib/docker" {
			return true, nil
		}
	}

	return false, nil
}

type DockerPartitionCheck struct {
	*CheckDefinitionImpl
	fstab string
}

func makeDockerPartitionCheck() Check {
	return &DockerPartitionCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:   "CIS-Docker-Benchmark-1.1",
			category:     "Host Configuration",
			name:         "Create a separate partition for containers",
			impact:       "None",
			rationale:    "Docker depends on /var/lib/docker as the default directory where all docker related files, including the images, are stored. This directory might fill up fast and soon Docker and the host could become unusable. So, it is advisable to create a separate partition (logical volume) for storing Docker files.",
			description:  `All Docker containers and their data and metadata is stored under /var/lib/docker directory. By default, /var/lib/docker would be mounted under / or /var partitions based on availability.`,
			defaultValue: `By default, /var/lib/docker would be mounted under / or /var partitions based on availability.`,
			references:   []string{"http://www.projectatomic.io/docs/docker-storage-recommendation"},
			remediation:  `For new installations, create a separate partition for /var/lib/docker mount point. For systems that were previously installed, use the Logical Volume Manager (LVM) to create partitions.`,
			auditDescription: `At the Docker host execute the below command:

grep /var/lib/docker /etc/fstab
		
This should return the partition details for /var/lib/docker mount point.`,
		},
		fstab: DEFAULT_FSTAB,
	}
}
