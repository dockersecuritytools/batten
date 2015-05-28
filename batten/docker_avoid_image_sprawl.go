package batten

import docker "github.com/fsouza/go-dockerclient"

func (dc *DockerAvoidImageSprawl) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerAvoidImageSprawl) AuditCheck() (bool, error) {
	client, err := docker.NewClient(DockerUnixSocket)

	if err != nil {
		// TODO: log error message
		return false, err
	}

	images, err := client.ListImages(docker.ListImagesOptions{All: false})

	if err != nil {
		// TODO: log error message
		return false, err
	}

	containers, err := client.ListContainers(docker.ListContainersOptions{All: false})

	if err != nil {
		// TODO: log error message
		return false, err
	}

	if len(images) > len(containers) {
		return false, nil
	}

	var uniqIds map[string]string = make(map[string]string, 0)

	for _, c := range containers {
		uniqIds[c.Image] = c.Image
	}

	if len(uniqIds) != len(images) {
		return false, nil
	}

	return true, nil
}

type DockerAvoidImageSprawl struct {
	*CheckDefinitionImpl
}

func makeDockerAvoidImageSprawl() Check {
	return &DockerAvoidImageSprawl{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-6.6",
			category:    "Docker Security Operations",
			name:        "Avoid image sprawl",
			impact:      "None",
			description: `Do not keep a large number of container images on the same host. Use only tagged images as appropriate.`,
			rationale:   `Tagged images are useful to fall back from "latest" to a specific version of an image in production. Images with unused or old tags may contain vulnerabilities that might be exploited, if instantiated. Additionally, if you fail to remove unused images from the system and there are various such redundant and unused images, the host filesystem may become full and could lead to denial of service.`,
			auditDescription: `Step 1 Make a list of all image IDs that are currently instantiated by executing below command: 
docker ps -q | xargs docker inspect --format '{{ .Id }}: Image={{ .Image }}' 

Step 2: List all the images present on the system by executing below command: 
docker images 

Step 3: Compare the list of image IDs populated from Step 1 and Step 2 and find out images that are currently not being instantiated. If any such unused or old images are found, discuss with the system administrator the need to keep such images on the system. If such a need is not justified enough, then this recommendation is non-compliant.`,
			remediation: `Keep the set of the images that you actually need and establish a workflow to remove old or stale images from the host. Additionally, use features such as pull-by-digest to get specific images from the registry. 
Additionally, you can follow below set of steps to find out unused images on the system and delete them. 
Step 1 Make a list of all image IDs that are currently instantiated by executing below command: 
docker ps -q | xargs docker inspect --format '{{ .Id }}: Image={{ .Image }}' 
Step 2: List all the images present on the system by executing below command: 
docker images 
Step 3: Compare the list of image IDs populated from Step 1 and Step 2 and find out images that are currently not being instantiated. 
Step 4: Decide if you want to keep the images that are not currently in use. If not delete them by executing below command: 
docker rmi $IMAGE_ID`,
			defaultValue: `Images and layered filesystems remain accessible on the host until the administrator removes all tags that refer to those images or layers.`,
			references: []string{
				"http://craiccomputing.blogspot.in/2014/09/clean-up-unused-docker-containers-and.html",
				"https://forums.docker.com/t/command-to-remove-all-unused-images/20/8",
			},
		},
	}
}
