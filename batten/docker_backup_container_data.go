package batten

func (dc *DockerBackupContainerData) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerBackupContainerData) AuditCheck() (bool, error) {
	// TODO
	return true, nil
}

type DockerBackupContainerData struct {
	*CheckDefinitionImpl
}

func makeDockerBackupContainerData() Check {
	return &DockerBackupContainerData{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:  "CIS-Docker-Benchmark-6.4",
			category:    "Docker Security Operations",
			name:        `Backup container data`,
			description: `Take regular backups of your container data volumes.`,
			rationale:   `Containers might run services that are critical for your business. Taking regular data backups would ensure that if there is ever any loss of data you would still have your data in backup. The loss of data could be devastating for your business.`,
			auditDescription: `Ask the system administrator whether container data volumes are regularly backed up. Verify a copy of the backup and ensure that the organization's backup policy is followed. 
Additionally, you can execute the below command for each container instance to list the changed files and directories in the container ̓s filesystem. Ideally, nothing should be stored on container's filesystem. 
docker diff $INSTANCE_ID`,
			remediation: `You should follow your organization's policy for data backup. You can take backup of your container data volume using '--volumes-from' parameter as below: 
$> docker run <Run arguments> --volumes-from $INSTANCE_ID -v [host-dir]:[container- dir] <Container Image Name or ID> <Command> 
For example, 
$> docker run --volumes-from 699ee3233b96 -v /mybackup:/backup centos tar cvf /backup/backup.tar /exampledatatobackup`,
			impact:       "None",
			defaultValue: `By default, no data backup happens for container data volumes.`,
			references: []string{
				"http://docs.docker.com/userguide/dockervolumes/#backup-restore-or-migrate-data-volumes",
			},
		},
	}
}
