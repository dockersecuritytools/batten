package batten

func (dc *DockerNoUnnecessaryPackagesCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerNoUnnecessaryPackagesCheck) AuditCheck() (bool, error) {
	// TODO: implement
	return true, nil
}

type DockerNoUnnecessaryPackagesCheck struct {
	*CheckDefinitionImpl
}

func makeDockerNoUnnecessaryPackagesCheck() Check {
	return &DockerNoUnnecessaryPackagesCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-4.3",
			category:    `Container Images and Build File`,
			name:        `Do not install unnecessary packages in the container`,
			description: `Containers tend to be minimal and slim down versions of the Operating System. Do not  install anything that does not justify the purpose of container. `,
			rationale:   `Bloating containers with unnecessary software could possibly increase the attack surface  of the container. This also voids the concept of minimal and slim down versions of  container images. Hence, do not install anything else apart from what is truly needed for  the purpose of the container. `,
			auditDescription: `Step 1: List all the running instances of containers by executing below command: 

docker ps -q

Step 2: For each container instance, execute the below or equivalent command: 

docker exec $INSTANCE_ID rpm -qa

The above command would list the packages installed on the container. Review the list and  ensure that it is legitimate. `,
			remediation: `At the outset, do not install anything on the container that does not justify the purpose. If  the image had some packages that your container does not use, uninstall them. 
Consider using a minimal base image rather than the standard Redhat/Centos/Debian  images if you can. Some of the options include BusyBox and Alpine.  
￼￼￼￼￼￼
Not only does this trim your image size from >150Mb to ~20 Mb, there are also fewer tools  and paths to escalate privileges. You can even remove the package installer as a final  hardening measure for leaf/production containers.`,
			impact:       `None`,
			defaultValue: `N/A`,
			references: []string{
				"https://docs.docker.com/userguide/dockerimages/",
				"http://www.livewyer.com/blog/2015/02/24/slimming-down-your-docker-containers-alpine-linux",
				"https://github.com/progrium/busybox",
			},
		},
	}
}
