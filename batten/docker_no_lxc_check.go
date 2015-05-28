package batten

import (
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

func (dc *DockerNoLxcCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

//
// AuditCheck looks for --exec-driver in the docker
// daemon options, e..g
//
// docker -d --exec-driver=lxc
//
func (dc *DockerNoLxcCheck) AuditCheck() (bool, error) {
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
	driver := info.Get("Execution Driver")

	if strings.Contains(driver, "lxc") {
		return false, nil
	}

	return true, nil

}

type DockerNoLxcCheck struct {
	*CheckDefinitionImpl
}

func makeDockerNoLxcCheck() Check {
	return &DockerNoLxcCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:   "CIS-Docker-Benchmark-2.1",
			category:     "Docker Daemon Configuration",
			name:         "Do not use lxc execution driver",
			description:  "The default Docker execution driver is 'libcontainer'. LXC as an execution driver is optional and just has legacy support.",
			rationale:    "There is still legacy support for the original LXC userspace tools via the 'lxc' execution driver, however, this is not where the primary development of new functionality is taking place. Docker out of the box can now manipulate namespaces, control groups, capabilities, apparmor profiles, network interfaces and firewalling rules - all in a consistent and predictable way, and without depending on LXC or any other userland package. This drastically reduces the number of moving parts, and insulates Docker from the side-effects introduced across versions and distributions of LXC.",
			defaultValue: `By default, Docker execution driver is 'libcontainer'`,
			auditDescription: `$> ps -ef | grep docker | grep lxc

The above command should not return anything. This would ensure that the '-e' or '-- exec-driver' parameter is either not present or not set to 'lxc'.`,
			references: []string{
				"http://www.infoq.com/news/2014/03/docker_0_9",
				"http://docs.docker.com/reference/commandline/cli/#docker-exec-driver-option",
			},
			remediation: `Do not run the Docker daemon with 'lxc' as execution driver. For example, do not start the Docker daemon as below:

$> docker -d --exec-driver=lxc`,
			impact: "None",
		},
	}
}
