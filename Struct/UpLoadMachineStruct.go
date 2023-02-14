package Struct

type UpLoadMachineStruct struct {
	Version UpLoadMachineVersionStruct `json:"version"`
	BinData string `json:"bindata"`
}

type UpLoadMachineVersionStruct struct {
	DnsVersion string `json:"dns_version"`
	DhcpVersion string `json:"dhcp_version"`
	AddVersion string `json:"add_version"`
	DdiVersion string `json:"ddi_version"`
	ChangePasswd string `json:"change_passwd"`
}