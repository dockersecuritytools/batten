package batten

import (
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
)

func (dc *DockerUseTrustedImagesCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerUseTrustedImagesCheck) AuditCheck() (bool, error) {
	client, err := getDockerAPIConnection()
	if err != nil {
		return false, err
	}

	images, err := client.ListImages(docker.ListImagesOptions{All: false})

	if err != nil {
		return false, err
	}

	// TODO: there needs to be a better way to check
	// this against a policy.
	for _, img := range images {

		if !stringInSlice(img.RepoTags[0], dc.trustedRepoTags) {
			imgDetails, err := client.InspectImage(img.ID)
			fmt.Println("XXX", imgDetails.Author, err)
			// TODO: log the failed repository
			// return false, nil
		}
	}

	return true, nil
}

type DockerUseTrustedImagesCheck struct {
	*CheckDefinitionImpl
	trustedRepoTags []string
}

func makeDockerUseTrustedImagesCheck() Check {
	return &DockerUseTrustedImagesCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-4.2",
			category:    `Container Images and Build File`,
			name:        `Use trusted base images for containers`,
			description: `Ensure that the container image is written either from scratch or is based on another established and trusted base image downloaded over a secure channel.`,
			rationale:   `Official repositories are Docker images curated and optimized by the Docker community or the vendor. But, the Docker container image signing and verification feature is not yet ready. Hence, the Docker engine does not verify the provenance of the container images by itself. You should thus exercise a great deal of caution when obtaining container images. `,
			auditDescription: `Inspect the Docker host by executing the below command:

$> docker images

This would list all the container images that are currently available for use on the Docker host. Interview the system administrator and obtain a proof of evidence that the list of images was obtained from trusted source over a secure channel.`,
			remediation:  `Only download the container images from a source you trust over a secure channel. Additionally, use features such as pull-by-digest to get specific images from the registry.`,
			impact:       `None`,
			defaultValue: `Not Applicable.`,
			references: []string{
				"https://titanous.com/posts/docker-insecurity",
				"https://registry.hub.docker.com/",
				"http://blog.docker.com/2014/10/docker-1-3-signed-images-process-injection-security-options-mac-shared-directories/",
				"https://github.com/docker/docker/issues/8093",
				"http://docs.docker.com/reference/commandline/cli/#pull",
				"https://github.com/docker/docker/pull/11109",
			},
		},
	}
}
