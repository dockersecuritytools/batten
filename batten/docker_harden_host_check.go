package batten

func (dc *DockerHardenHostCheck) GetCheckDefinition() CheckDefinition {
	return dc
}

func (dc *DockerHardenHostCheck) AuditCheck() (bool, error) {
	if dc.policy != "" {
		// TODO: run the policy check command
	}
	return true, nil
}

type DockerHardenHostCheck struct {
	*CheckDefinitionImpl
	policy string
}

func makeDockerHardenHostCheck() Check {
	return &DockerHardenHostCheck{
		CheckDefinitionImpl: &CheckDefinitionImpl{
			identifier:       "CIS-Docker-Benchmark-1.4",
			name:             "Harden the container host",
			impact:           "None",
			description:      "Containers run on a Linux host. A container host can run one or more containers. It is of utmost importance to harden the host to mitigate host security misconfiguration.",
			rationale:        "You should follow infrastructure security best practices and harden your host OS. Keeping the host system hardened would ensure that the host vulnerabilities are mitigated. Not hardening the host system could lead to security exposures and breaches.",
			defaultValue:     "By default, host has factory settings. It is not hardened.",
			auditDescription: "Ensure that the host specific security guidelines are followed. Ask the system administrators which security benchmark does current host system comply with. Ensure that the host systems actually comply with that host specific security benchmark.",
			references: []string{
				"https://docs.docker.com/articles/security/",
				"https://benchmarks.cisecurity.org/downloads/multiform/index.cfm",
				"http://docs.docker.com/articles/security/#other-kernel-security-features 4. http://blog.dotcloud.com/kernel-secrets-from-the-paas-garage-part-44-g 5. https://grsecurity.net/",
				"https://en.wikibooks.org/wiki/Grsecurity",
				"https://pax.grsecurity.net/",
				"http://en.wikipedia.org/wiki/PaX",
			},
			remediation: "You may consider various CIS Security Benchmarks for your container host. If you have other security guidelines or regulatory requirements to adhere to, please follow them as suitable in your environment.  Additionally, you can run a kernel with grsecurity and PaX. This would add many safety checks, both at compile-time and run-time. It is also designed to defeat many exploits and has powerful security features. These features do not require Docker-specific configuration, since those security features apply system-wide, independent of containers.",
		},
	}
}
