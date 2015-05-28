package batten

type CheckDefinitionImpl struct {
	name             string
	category         string
	description      string
	rationale        string
	remediation      string
	impact           string
	defaultValue     string
	references       []string
	auditDescription string
	identifier       string
}

func (c *CheckDefinitionImpl) Category() string {
	return c.category
}

func (c *CheckDefinitionImpl) AuditDescription() string {
	return c.auditDescription
}

func (c *CheckDefinitionImpl) Identifier() string {
	return c.identifier
}

func (c *CheckDefinitionImpl) Name() string {
	return c.name
}

func (c *CheckDefinitionImpl) Description() string {
	return c.description
}

func (c *CheckDefinitionImpl) Rationale() string {
	return c.rationale
}

func (c *CheckDefinitionImpl) Remediation() string {
	return c.remediation
}
func (c *CheckDefinitionImpl) Impact() string {
	return c.impact
}

func (c *CheckDefinitionImpl) DefaultValue() string {
	return c.defaultValue
}

func (c *CheckDefinitionImpl) References() []string {
	return c.references
}

type CheckDefinition interface {
	Identifier() string
	Name() string
	Description() string
	Rationale() string
	Remediation() string
	Impact() string
	DefaultValue() string
	References() []string
}

type Check interface {
	AuditCheck() (bool, error)
	GetCheckDefinition() CheckDefinition
}

type CheckResults struct {
	Success         bool
	Error           error
	CheckDefinition CheckDefinition
}

func RunCheck(c Check) *CheckResults {
	succ, err := c.AuditCheck()

	return &CheckResults{
		Success:         succ,
		Error:           err,
		CheckDefinition: c.GetCheckDefinition(),
	}
}

// TODO: put the checks in a diff package and
// allow to register with the batten main package.
var Checks []Check = []Check{
	// 1.x host checks
	makeDockerPartitionCheck(),
	makeDockerKernelCheck(),
	makeDockerDevToolsCheck(),
	makeDockerHardenHostCheck(),
	makeDockerRemoveNonEssentialSvcsCheck(),
	makeDockerVersionCheck(),
	makeDockerTrustedUsersCheck(),
	makeDockerDaemonAuditingCheck(),
	makeDockerAuditFilesVarLibDocker(),
	makeDockerAuditFilesEtcDocker(),
	makeDockerAuditFilesDockerRegistry(),
	makeDockerAuditFilesDockerService(),
	makeDockerAuditFilesDockerSock(),
	makeDockerAuditFilesSysconfigDocker(),
	makeDockerAuditFilesSysconfigDockerNetwork(),
	makeDockerAuditFilesSysconfigDockerRegistry(),
	makeDockerAuditFilesSysconfigDockerStorage(),
	makeDockerAuditFilesEtcDefaultDocker(),
	// 2.x Docker daemon configuration checks
	makeDockerNoLxcCheck(),
	makeDockerRestrictedNetworkTrafficCheck(),
	makeDockerSetLoggingLevelCheck(),
	makeDockerEnableIptablesCheck(),
	makeDockerInsecureRegistriesCheck(),
	makeDockerLocalRegistryCheck(),
	makeDockerNoAufsCheck(),
	makeDockerPortCheck(),
	makeDockerTLSCheck(),
	makeDockerUlimitCheck(),
	// 3.x Docker daemon configuration files
	makeDockerSvcOwnerCheck(),
	makeDockerSvcFilePermsCheck(),
	makeDockerRegistrySvcOwnerCheck(),
	makeDockerRegistrySvcFilePermsCheck(),
	makeDockerSystemdSocketOwnerCheck(),
	makeDockerSystemdSocketFilePermsCheck(),
	makeDockerEnvFileOwnerCheck(),
	makeDockerEnvFilePermsCheck(),
	makeDockerNetworkEnvOwnerCheck(),
	makeDockerNetworkEnvFilePermsCheck(),
	makeDockerRegistryEnvOwnerCheck(),
	makeDockerRegistryEnvFilePermsCheck(),
	makeDockerStorageEnvOwnerCheck(),
	makeDockerStorageEnvFilePermsCheck(),
	makeDockerEtcDockerOwnerCheck(),
	makeDockerEtcDockerFilePermsCheck(),
	makeDockerRegistryCertsOwnerCheck(),
	makeDockerRegistryCertsFilePermsCheck(),
	makeDockerTLSCACertOwnerCheck(),
	makeDockerTLSCACertFilePermsCheck(),
	makeDockerTLSCertOwnerCheck(),
	makeDockerTLSCertFilePermsCheck(),
	makeDockerTLSKeyOwnerCheck(),
	makeDockerTLSKeyFilePermsCheck(),
	makeDockerSocketOwnerCheck(),
	makeDockerSocketFilePermsCheck(),

	// 4.x
	makeDockerContainerUserCheck(),
	makeDockerUseTrustedImagesCheck(),
	makeDockerNoUnnecessaryPackagesCheck(),

	// 5.x
	makeDockerVerifyAppArmorProfile(),
	makeDockerVerifySELinuxProfile(),
	makeDockerSingleMainProcess(),
	makeDockerRestrictKernel(),

	// 6.x
	makeDockerPerformSecurityAudits(),
	makeDockerMonitorContainers(),
	makeDockerCheckEndpointProtectionPlatform(),
	makeDockerBackupContainerData(),
	makeDockerCheckCentralLogCollection(),
	makeDockerAvoidImageSprawl(),
	makeDockerAvoidContainerSprawl(),
}
