package batten

import (
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

func (dc *DockerNoAufsCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerNoAufsCheck) AuditCheck() (bool, error) {
	client, err := docker.NewClient(DockerUnixSocket)

	if err != nil {
		// TODO: log error message
		return false, err
	}

	info, err := client.Info()

	if err != nil {
		// TODO: log error message
		return false, err
	}
	driver := info.Get("Driver")
	if strings.Contains(driver, "aufs") {
		return false, nil
	}

	return true, nil
}

type DockerNoAufsCheck struct {
	*CheckDefinitionImpl
}

func makeDockerNoAufsCheck() Check {
	return &DockerNoAufsCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-2.7",
			category:    `Docker daemon configuration`,
			name:        `Do not use the aufs storage driver`,
			description: `Do not use 'aufs' as storage driver for your Docker instance.`,
			rationale:   `The 'aufs' storage driver is the oldest storage driver. It is based on a Linux kernel patch-set that is unlikely to be merged into the main Linux kernel. 'aufs' driver is also known to cause some serious kernel crashes. 'aufs' just has legacy support from Docker. Most importantly, 'aufs' is not a supported driver in many Linux distributions using latest Linux kernels.`,
			auditDescription: `Execute the below command and verify that 'aufs' is not used as storage driver: 
			
docker info | grep -e "^Storage Driver:\s*aufs\s*$"

The above command should not return anything.`,
			remediation: `Do not explicitly use 'aufs' as storage driver.

For example, do not start Docker daemon as below:
$> docker -s aufs -d`,
			impact:       `'aufs' is the only storage driver that allows containers to share executable and shared library memory. It might be useful if you are running thousands of containers with the same program or libraries.`,
			defaultValue: `By default, Docker uses 'devicemapper' as the storage driver on most of the platforms. Default storage driver can vary based on your OS vendor. You should use the storage driver that is best supported by your preferred vendor.`,
			references: []string{

				"http://docs.docker.com/reference/commandline/cli/#daemon-storage-driver-option",
				"https://github.com/docker/docker/issues/6047",
				"http://muehe.org/posts/switching-docker-from-aufs-to-devicemapper/",
				"http://jpetazzo.github.io/assets/2015-03-05-deep-dive-into-docker-storage-drivers.html#1",
			},
		},
	}
}
