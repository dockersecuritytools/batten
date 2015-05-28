package batten

import "github.com/fsouza/go-dockerclient"

func (dc *DockerRestrictKernel) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerRestrictKernel) AuditCheck() (bool, error) {
	client, err := docker.NewClient(DockerUnixSocket)

	if err != nil {
		// TODO: log error message
		return false, err
	}

	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})

	if err != nil {
		// TODO: log error message
		return false, err
	}

	for _, c := range containers {
		if cc, err := client.InspectContainer(c.ID); err == nil {
			if len(cc.HostConfig.CapAdd) <= 0 {
				return true, nil
			}

			// TODO: do better check here.
			for _, sysc := range cc.HostConfig.CapAdd {
				if dc.blockedCalls[sysc] {
					return false, nil
				}
			}
		}
	}

	return true, nil
}

type DockerRestrictKernel struct {
	*CheckDefinitionImpl
	blockedCalls map[string]bool
}

func makeDockerRestrictKernel() Check {
	return &DockerRestrictKernel{
		blockedCalls: map[string]bool{
			"NET_ADMIN":  true,
			"SYS_ADMIN":  true,
			"SYS_MODULE": true,
		},
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-5.4",
			category:    `Container Runtime`,
			name:        `Restrict Linux Kernel Capabilities within containers`,
			description: `By default, Docker starts containers with a restricted set of Linux Kernel Capabilities. It means that any process may be granted the required capabilities instead of root access. Using Linux Kernel Capabilities, the processes do not have to run as root for almost all the specific areas where root privileges are usually needed.`,
			rationale: `Docker supports the addition and removal of capabilities, allowing use of a non-default profile. This may make Docker more secure through capability removal, or less secure through the addition of capabilities. It is thus recommended to remove all capabilities except those explicitly required for your container process. 
			For example, capabilities such as below are usually not needed for container process: 
     
NET_ADMIN
SYS_ADMIN
SYS_MODULE
`,
			auditDescription: `docker ps -q | xargs docker inspect --format '{{ .Id }}: CapAdd={{ .HostConfig.CapAdd }} CapDrop={{ .HostConfig.CapDrop }}' 
Verify that the added and dropped Linux Kernel Capabilities are in line with the ones needed for container process for each container instance.`,
			remediation: `Execute the below command to add needed capabilities: 
$> docker run --cap-add={"Capability 1","Capability 2"} <Run arguments> <Container Image Name or ID> <Command> 
For example, 

docker ps -q | xargs docker inspect --format '{{ .Id }}: CapAdd={{ .HostConfig.CapAdd }} CapDrop={{ .HostConfig.CapDrop }}' 
          
$> docker run --cap-add={"NET_ADMIN","SYS_ADMIN"} -i -t centos:latest /bin/bash 

$> docker run --cap-drop={"Capability 1","Capability 2"} <Run arguments> <Container Image Name or ID> <Command> 
For example, 
$> docker run --cap-drop={"SETUID","SETGID"} -i -t centos:latest /bin/bash `,
			impact: `Based on what Linux Kernel Capabilities were added or dropped, restrictions within the container would apply.`,
			defaultValue: `By default, below capabilities are available for containers: 
AUDIT_WRITE 
CHOWN 
DAC_OVERRIDE 
FOWNER 
FSETID 
KILL MKNOD 
NET_BIND_SERVICE 
NET_RAW 
SETFCAP 
SETGID 
SETPCAP 
SETUID 
SYS_CHROOT`,
			references: []string{
				"https://docs.docker.com/articles/security/#linux-kernel-capabilities",
				"https://github.com/docker/docker/blob/master/daemon/execdriver/native/temp/late/default_template.go",
				"http://man7.org/linux/man-pages/man7/capabilities.7.html",
			},
		},
	}
}
