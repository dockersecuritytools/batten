package batten

import (
	"errors"
	"os/exec"
	"strings"

	version "github.com/hashicorp/go-version"
)

func (dc *DockerKernelCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerKernelCheck) AuditCheck() (bool, error) {

	cmd := exec.Command("uname", "-r")

	bytes, err := cmd.CombinedOutput()

	if err != nil {
		return false, err
	}

	lines := strings.Split(string(bytes), "\n")

	if len(lines) < 1 {
		return false, errors.New("Nothing returned from uname -r")
	}

	kernelstring := lines[0]
	parts := strings.Split(kernelstring, "-")

	if len(parts) < 1 {
		return false, errors.New("Malformed kernel string" + kernelstring)
	}
	kernelversion := parts[0]

	v1, err := version.NewVersion(kernelversion)
	if err != nil {
		return false, err
	}
	targetVersion, err := version.NewVersion("3.10")
	if err != nil {
		return false, err
	}

	if v1.Compare(targetVersion) >= 0 {
		return true, nil
	}
	// TODO: print out kernel version
	return false, nil
}

type DockerKernelCheck struct {
	*CheckDefinitionImpl
}

func makeDockerKernelCheck() Check {
	return &DockerKernelCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:   "CIS-Docker-Benchmark-1.2",
			category:     "Host Configuration",
			name:         "Use the updated Linux Kernel",
			impact:       "None",
			description:  `Docker in daemon mode has specific kernel requirements. A 3.10 Linux kernel is the minimum requirement for Docker.`,
			rationale:    `Kernels older than 3.10 lack some of the features required to run Docker containers. These older versions are known to have bugs which cause data loss and frequently panic under certain conditions. The latest minor version (3.x.y) of the 3.10 (or a newer maintained version) Linux kernel is thus recommended. Additionally, using the updated Linux kernels ensures that critical kernel bugs found earlier are fixed.`,
			defaultValue: `N/A`,
			auditDescription: `Execute the below command to find out Linux kernel version:

uname -r

Ensure that the kernel version found is 3.10 or newer.`,
			references:  []string{"https://docs.docker.com/installation/binaries/#check-kernel-dependencies", "https://docs.docker.com/installation/#installation-list"},
			remediation: `Check out the Docker kernel and OS requirements and suitably choose your kernel and OS.`,
		},
	}
}
