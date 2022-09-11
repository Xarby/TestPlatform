package Struct

// 部署节点的结构

// FTP部署的结构
type FtpTask struct {
	Ftp FtpStruct `json:"ftp"`
	DnsVersion string `json:"dns_version"`
	AddVersion string `json:"add_version"`
	DhcpVersion string `json:"dhcp_version"`
	ZddiPath string `json:"zddi_path"`
	BuildPath string `json:"build_path"`
	ZddiDevices []ScpStruct `json:"zddi_devices"`
}
//  SCP部署的结构
type ScpTask struct {
	GetScpZddi ScpStruct `json:"get_scp_zddi"`
	GetScpBuild ScpStruct  `json:"get_scp_build"`
	DnsVersion string `json:"dns_version"`
	AddVersion string `json:"add_version"`
	DhcpVersion string `json:"dhcp_version"`
	ZddiDevices []ScpStruct `json:"zddi_devices"`
}



