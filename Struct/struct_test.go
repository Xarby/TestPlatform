package Struct

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDegDomain(t *testing.T) {
	ftp_dev := FtpTask{}
	json.Unmarshal([]byte(`{
	"ftp": {
		"ipaddr": "10.1.101.37",
		"port": "21",
		"user": "zdns",
		"password": "zdns"
	},
	"dns_version": "3",
	"dhcp_version": "1",
	"add_version": "2",
	"zddi_path": "/3.15/3.15.2.3/zddi-3.15.2.3-20220809.el7.x86_64.rpm",
	"build_path": "/3.15/3.15.2.3/zddi_build-3.15.3.7-20220706.el7.x86_64.rpm",
	"zddi_devices": [
		{
			"ipaddr": "10.1.107.50",
			"user": "root",
			"port": "22",
			"role": "master",
			"password": "zdns@knet.cn"
		}
	]
}`), &ftp_dev)
	fmt.Println(ftp_dev.CheckFtpTask())
}
