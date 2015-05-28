package batten

import "github.com/fsouza/go-dockerclient"

func (dc *DockerVerifyAppArmorProfile) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerVerifyAppArmorProfile) AuditCheck() (bool, error) {
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

	var count int
	for _, c := range containers {

		if cc, err := client.InspectContainer(c.ID); err == nil {
			if len(cc.AppArmorProfile) > 0 {
				count++
			}
		}
	}

	if count != len(containers) {
		return false, nil
	}

	return true, nil
}

type DockerVerifyAppArmorProfile struct {
	*CheckDefinitionImpl
}

func makeDockerVerifyAppArmorProfile() Check {
	return &DockerVerifyAppArmorProfile{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:       "CIS-Docker-Benchmark-5.1",
			category:         "Container Runtime",
			name:             `Verify AppArmor Profile, if applicable`,
			description:      `AppArmor is an effective and easy-to-use Linux application security system. It is available on quite a few Linux distributions by default such as Debian and Ubuntu.`,
			rationale:        `AppArmor protects the Linux OS and applications from various threats by enforcing security policy which is also known as AppArmor profile. You should create a AppArmor profile for your containers. This would enforce security policies on the containers as defined in the profile.`,
			auditDescription: `The above command should return a valid AppArmor Profile for each container instance.`,
			remediation: `If AppArmor is applicable for your Linux OS, use it. You may have to follow below set of steps:
1.	Verify if AppArmor is installed. If not, install it.  
2.	Create or import a AppArmor profile for Docker containers.  
3.	Put this profile in enforcing mode.  
4.	Start your Docker container using the Docker AppArmor profile. For example,  docker run -i -t --security-opt="apparmor:PROFILENAME" centos /bin/bash  
     
docker ps -q | xargs docker inspect --format '{{ .Id }}: AppArmorProfile={{ .AppArmorProfile }}'`,
			impact:       "The container (process) would have set of restrictions as defined in AppArmor profile. If your AppArmor profile is mis-configured, then the container may not entirely work as expected.",
			defaultValue: `By default, no AppArmor profiles are applied on containers.`,
			references: []string{
				"http://docs.docker.com/articles/security/#other-kernel-security-features",
				"http://docs.docker.com/reference/run/#security-configuration",
				"http://wiki.apparmor.net/index.php/Main_Page",
			},
		},
	}
}
