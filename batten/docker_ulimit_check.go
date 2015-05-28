package batten

import "errors"

import "github.com/jandre/procfs/limits"

func (dc *DockerUlimitCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerUlimitCheck) checkForProcUlimits(l *limits.Limits) bool {
	hardVal := l.Processes.HardValue
	if dc.processesUlimitMinimum != 0 {
		return hardVal >= dc.processesUlimitMinimum
	} else {
		return hardVal != limits.UNLIMITED
		// just make sure it's not "unlimited"
	}
	return true
}

// TODO: make this testable
func (dc *DockerUlimitCheck) checkForFileUlimits(l *limits.Limits) bool {

	hardVal := l.OpenFiles.HardValue
	if dc.openFilesUlimitMinimum != 0 {

		return hardVal >= dc.openFilesUlimitMinimum
	} else {
		return hardVal != limits.UNLIMITED
		// just make sure it's not "unlimited"
	}

	return true
}

func (dc *DockerUlimitCheck) AuditCheck() (bool, error) {

	process, err := getDockerProcess(dc.dockerPidFile)

	if err != nil {
		return false, err
	}

	if process != nil {
		l, err := process.Limits()
		if err != nil {
			return false, err
		}
		if dc.checkForFileUlimits(l) && dc.checkForProcUlimits(l) {
			return true, nil
		} else {
			return false, nil
		}
	}

	return false, errors.New("Docker daemon not running")

}

type DockerUlimitCheck struct {
	*CheckDefinitionImpl
	// TODO: make configurable
	dockerPidFile          string
	processesUlimitMinimum int
	openFilesUlimitMinimum int
}

func makeDockerUlimitCheck() Check {
	return &DockerUlimitCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-2.10",
			category:    `Docker daemon configuration`,
			name:        `Set default ulimit as appropriate`,
			description: `Set the default ulimit options as appropriate in your environment.`,
			rationale:   `ulimit provides control over the resources available to the shell and to processes started by it. Setting system resource limits judiciously saves you from many disasters such as a fork bomb. Sometimes, even friendly users and legitimate processes can overuse system resources and in-turn can make the system unusable.  Setting default ulimit for the Docker daemon would enforce the ulimit for all container instances. You would not need to setup ulimit for each container instance. However, the default ulimit can be overridden during container runtime, if needed. Hence, to control the system resources, define a default ulimit as needed in your environment.`,
			auditDescription: `$> ps -ef | grep docker

Ensure that the '--default-ulimit' parameter is set as appropriate.`,
			remediation: `Run the docker in daemon mode and pass '--default-ulimit' as argument with respective ulimits as appropriate in your environment.
			
For Example,

$> docker -d --default-ulimit nproc=1024:2408 --default-ulimit nofile=100:200`,
			impact:       `If the ulimits are not set properly, the desired resource control might not be achieved and might even make the system unusable.`,
			defaultValue: `By default, no ulimit is set.`,
			references: []string{
				"http://docs.docker.com/reference/commandline/cli/#default-ulimits",
			},
		},
		processesUlimitMinimum: 0,
		openFilesUlimitMinimum: 0,
	}
}
