package batten

import "github.com/fsouza/go-dockerclient"

func (dc *DockerVerifySELinuxProfile) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerVerifySELinuxProfile) AuditCheck() (bool, error) {
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
			if len(cc.HostConfig.SecurityOpt) > 0 {
				count++
			}
		}
	}

	if count != len(containers) {
		return false, nil
	}

	return true, nil
}

type DockerVerifySELinuxProfile struct {
	*CheckDefinitionImpl
}

func makeDockerVerifySELinuxProfile() Check {
	return &DockerVerifySELinuxProfile{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-5.2",
			category:    "Container Runtime",
			name:        `Verify SELinux security options, if applicable`,
			description: `SELinux is an effective and easy-to-use Linux application security system. It is available on quite a few Linux distributions by default such as Red Hat and Fedora.`,
			rationale:   `SELinux provides a Mandatory Access Control (MAC) system that greatly augments the default Discretionary Access Control (DAC) model. You can thus add an extra layer of safety by enabling SELinux on your Linux host, if applicable.`,
			auditDescription: `The above command should return all the security options currently configured for the containers.
docker ps -q | xargs docker inspect --format '{{ .Id }}: SecurityOpt={{ .HostConfig.SecurityOpt }}'`,
			remediation: `If SELinux is applicable for your Linux OS, use it. You may have to follow below set of steps: 
1.  Set the SELinux State. 
2.	Set the SELinux Policy.  
3.	Create or import a SELinux policy template for Docker containers.  
4.	Start Docker in daemon mode with SELinux enabled. For example,  docker -d --selinux-enabled  
5.	Start your Docker container using the security options. For example,  docker run -i -t --security-opt label:level:TopSecret centos /bin/bash`,
			impact:       `The container (process) would have set of restrictions as defined in SELinux policy. If your SELinux policy is mis-configured, then the container may not entirely work as expected.`,
			defaultValue: `By default, no SELinux security options are applied on containers.`,
			references: []string{
				"http://docs.docker.com/articles/security/#other-kernel-security-features",
				"http://docs.docker.com/reference/run/#security-configuration",
				"http://docs.fedoraproject.org/en-US/Fedora/13/html/Security-Enhanced_Linux/",
			},
		},
	}
}
