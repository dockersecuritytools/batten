package batten

import (
	"errors"
	"strings"
)

func (dc *DockerLocalRegistryCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerLocalRegistryCheck) lookForRegistry(argv string) (bool, error) {

	if strings.Contains(argv, "--registry-mirror=") {
		return true, nil
	}

	return false, nil
}
func (dc *DockerLocalRegistryCheck) AuditCheck() (bool, error) {

	succ, args, err := readDockerDaemonArgs(dc.dockerPidFile)

	if err != nil {
		return false, err
	}

	if succ {
		argv := strings.Join(args, " ")
		return dc.lookForRegistry(argv)
	}
	return false, errors.New("Docker daemon not running")
}

type DockerLocalRegistryCheck struct {
	*CheckDefinitionImpl
	// TODO: make configurable
	dockerPidFile string
}

func makeDockerLocalRegistryCheck() Check {
	return &DockerLocalRegistryCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-2.6",
			category:    `Docker daemon configuration`,
			name:        `Setup a local registry mirror`,
			description: `The local registry mirror is serves the images from its own storage.`,
			rationale:   `If you have multiple instances of Docker running in your environment, each time one of them requires an image, it will have to go out to the internet and fetch it from public or your private Docker registry. By running a local registry mirror, you can keep image fetch traffic on your local network. So, your Docker instances need not have to be internet facing and thus this drastically reduces the threat vector. Additionally, it allows you to manage and securely store your images within your own environment.`,
			auditDescription: `Execute the below command to find out if a local registry is used:
 ps -ef | grep docker
			￼
Ensure that the '--registry-mirror' parameter is present.`,
			remediation: `Configure a local registry mirror and then start the Docker daemon as below:

$> docker --registry-mirror=<registry path> -d

For example,

$> docker --registry-mirror=https://10.0.0.2:5000 -d`,
			impact:       `The local registry mirror would need to be managed. It must have verified images that you use in your environment and those images must be kept updated time to time.`,
			defaultValue: `By default, there are no local registry mirrors setup on the host with Docker installation.`,
			references: []string{
				"http://docs.docker.com/articles/registry_mirror/",
			},
		},

		dockerPidFile: "/var/run/docker.pid",
	}
}
