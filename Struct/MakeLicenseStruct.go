package Struct

type MakeLicenseStruct struct {
	SshStruct SshStruct `json:"ssh_struct"`
	DnsVersion string `json:"dns_version"`
	DhcpVersion string `json:"dhcp_version"`
	AddVersion string `json:"add_version"`
	DdiVersion string `json:"ddi_version"`
}
