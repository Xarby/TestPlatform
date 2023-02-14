package Util

import (
	"TestPlatform/Const"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestTryConn(t *testing.T) {
	//fmt.Println("-----------")
	//fmt.Println(TryConn("10.1.104.55",443))
	//body := []byte(`{
	//"name": "slave",
	//"ip": "10.1.107.51",
	//"username": "admin",
	//"password": "admincns",
	//"role": "slave",
	//"group": "local",
	//"is_extend":"no"}`)
	//fmt.Println(string(body))
	//fmt.Println(PostRequests("https://10.1.107.54:20120/groups/local/members", body))
	//fmt.Println("-----------")

	//out, err := exec.Command("/bin/bash", "-c", "license_new -p '/root/TestPlatform/File/Private/tempLicense/10_1_121_28/pub.key' -q '/root/TestPlatform/File/Private/tempLicense/10_1_121_28/pri.key' -i '/root/TestPlatform/File/Private/tempLicense/10_1_121_28/machine.info' -l '/root/TestPlatform/File/Private/tempLicense/10_1_121_28/license.file' -r 'license_start_time:#:2022-12-03 00:00:00;license_valid_days:#:-1;license_remote_auth_flag:#:0;root_change_password_flag:#:0;device_type:#:0' -e 'DNSv3,DHCPv1,ADDv2 0 0'").CombinedOutput()
	//fmt.Println("license_new -p '/root/TestPlatform/File/Private/tempLicense/10_1_121_28/pub.key' -q '/root/TestPlatform/File/Private/tempLicense/10_1_121_28/pri.key' -i '/root/TestPlatform/File/Private/tempLicense/10_1_121_28/machine.info' -l '/root/TestPlatform/File/Private/tempLicense/10_1_121_28/license.file' -r 'license_start_time:#:2022-12-03 00:00:00;license_valid_days:#:-1;license_remote_auth_flag:#:0;root_change_password_flag:#:0;device_type:#:0' -e 'DNSv3,DHCPv1,ADDv2 0 0'")
	//if err != nil {
	//	fmt.Println(string(out), err.Error())
	//}else {
	//	fmt.Println(string(out), err.Error())
	//}
	programPath := "/root/TestPlatform/"

	tmpFilePath := programPath+Const.TempLicensePath+strings.Replace("10.1.121.111", ".", "_", -1)+"/"
	local_machine_info:=tmpFilePath+"machine.info"
	local_pub_key:=tmpFilePath+"pub.key"
	local_pri_key:=tmpFilePath+"pri.key"
	local_licnese_file:=tmpFilePath+"license.file"


	//fmt.Println(fmt.Sprintf("license_new -p '%s' -q '%s' -i '%s' -l '%s' -r 'license_start_time:#:2022-12-03 00:00:00;license_valid_days:#:-1;license_remote_auth_flag:#:0;root_change_password_flag:#:0;device_type:#:0' -e 'DNSv%s,DHCPv%s,ADDv%s 0 0'",local_pub_key,local_pri_key,local_machine_info,local_licnese_file, strconv.Itoa(3), strconv.Itoa(1), strconv.Itoa(2)))
	//OsExecCmd(fmt.Sprintf("license_new -p '%s' -q '%s' -i '%s' -l '%s' -r 'license_start_time:#:2022-12-03 00:00:00;license_valid_days:#:-1;license_remote_auth_flag:#:0;root_change_password_flag:#:0;device_type:#:0' -e 'DNSv%s,DHCPv%s,ADDv%s 0 0'",local_pub_key,local_pri_key,local_machine_info,local_licnese_file, strconv.Itoa(3), strconv.Itoa(1), strconv.Itoa(2)))

	cmds := fmt.Sprintf("license_old -p '%s' -q '%s' -i '%s' -l '%s' -r 'license_start_time:#:2022-12-03 13:22:40;license_valid_days:#:-1;license_remote_auth_flag:#:0;license_flag:#:1;root_change_password_flag:#:0;license_control_node_num:#:1;license_normal_flag:#:1;license_virtual_validate:#:1;user_name:#:test;license_type:#:0;device_type:#:0;elk:#:1;sgcloud_st_time:#:2022-12-03 13:22:40;sgcloud_vaild_days:#:-1;sgcloud_role:#:1' -e 'DNSv%s,DHCPv%s,ADDv%s,REGv0 0 0'", local_pub_key, local_pri_key, local_machine_info, local_licnese_file, strconv.Itoa(3), strconv.Itoa(1), strconv.Itoa(2))
	fmt.Println(cmds)
	OsExecCmd(cmds)
}