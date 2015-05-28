package batten

func (dc *DockerSecurityPatchesCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerSecurityPatchesCheck) AuditCheck() (bool, error) {
	// TODO: implement
	return true, nil
}

type DockerSecurityPatchesCheck struct {
	*CheckDefinitionImpl
}

func makeDockerSecurityPatchesCheck() Check {
	return &DockerSecurityPatchesCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-4.4",
			category:    `Container Images and Build File`,
			name:        `Rebuild the images to include security patches`,
			description: `Instead of patching your containers and images, rebuild the images from scratch and instantiate new containers from it.`,
			rationale: `Security patches are updates to products to resolve known issues. These patches update  the system to the most recent code base. Being on the current code base is important  because that's where vendors focus on fixing problems. Evaluate the security patches  before applying and follow the patching best practices. 
On conventional systems, rebuilding the system is a risky operation, because it is rarely  done, and many components can diverge over time between two rebuilds. However, when  using state-of-the-art configuration management or a build system like the one provided by  Docker with Dockerfiles, rebuilds are easy to do, and fast. This allows (and promotes)  frequent rebuilds, and ensures that at any given time, it is safe (and reliable) to do such a  full rebuilds. As such, in case of security vulnerability affecting a package or library,  rebuilding with the latest available software is the best practice and should be followed. `,
			auditDescription: `Step 1: List all the running instances of containers by executing below command: 
docker ps -q

Step 2: For each container instance, execute the below or equivalent command to find the  list of packages installed within container. Ensure that the security updates for various  affected packages are installed. 

docker exec INSTANCE_ID rpm -qa`,
			remediation: `Follow the below steps to rebuild the images with security patches: 
 
Step 1: 'docker pull' all the base images (i.e., given your set of Dockerfiles, extract all  images declared in 'FROM' instructions, and re-pull them to check for an updated version).  Step 2: Force a rebuild of each image with 'docker build --no-cache'.  Step 3: Restart all containers with the updated images. `,
			impact: `Rebuilding the images has to be done after upstream packages are available, otherwise re- pulling and rebuilding will do no good. When the affected packages are in the base image, it  is necessary to pull it (and therefore rebuild). When the affected packages are in the  downloaded packages, it is not necessary to pull the image; but nonetheless, in doubt, it is  recommended to always follow this strict procedure and rebuild the entire image. 
Note: If updated packages are not available and it is critical to install a security patch, live  patching could be used. `,
			defaultValue: `By default, containers and images are not updated of their own. `,
			references: []string{
				"https://docs.docker.com/userguide/dockerimages/",
			},
		},
	}
}
