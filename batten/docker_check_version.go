package batten

import (
	"errors"

	docker "github.com/fsouza/go-dockerclient"
	version "github.com/hashicorp/go-version"
)

func (dc *DockerVersionCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerVersionCheck) AuditCheck() (bool, error) {

	client, err := docker.NewClient(DockerUnixSocket)

	if err != nil {
		// TODO: log error message
		return false, err
	}

	v, err := client.Version()

	if err != nil {
		// TODO: log error message
		return false, err
	}

	if v == nil {
		return false, errors.New("Unable to retrieve docker version from api")
	}

	dockerVersion, err := version.NewVersion(v.Get("Version"))
	if err != nil {
		return false, err
	}
	targetVersion, err := version.NewVersion(dc.targetVersion)
	if err != nil {
		return false, err
	}

	if dockerVersion.Compare(targetVersion) >= 0 {
		return true, nil
	}
	return true, nil
}

type DockerVersionCheck struct {
	*CheckDefinitionImpl
	endpoint      string
	targetVersion string
}

func makeDockerVersionCheck() Check {
	return &DockerVersionCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier: "CIS-Docker-Benchmark-1.6",
			category:   "Host Configuration",
			name:       "Keep Docker up to date",
			impact:     "None",
			description: `The docker container solution is evolving to maturity and stability at a rapid pace. Like any other software, the vendor releases regular updates for Docker software that address security vulnerabilities, product bugs and bring in new functionality.

Keep a tab on these product updates and upgrade as frequently as when new security vulnerabilities are fixed.`,
			rationale:    "By staying up to date on Docker updates, vulnerabilities in the Docker software can be mitigated. An educated attacker may exploit known vulnerabilities when attempting to attain access or elevate privileges. Not installing regular Docker updates may leave you with running vulnerable Docker software. It might lead to elevation privileges, unauthorized access or other security breaches.",
			defaultValue: `N/A`,
			auditDescription: `Execute the below command and verify that the Docker version is up to date.

docker version`,
			references:  []string{"https://docs.docker.com/installation/"},
			remediation: "Download and install the updated Docker software from official Docker repository.",
		},

		targetVersion: "1.6",
		endpoint:      "unix:///var/run/docker.sock",
	}
}
