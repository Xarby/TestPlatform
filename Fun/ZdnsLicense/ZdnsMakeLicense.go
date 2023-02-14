package ZdnsLicense

import (
	"TestPlatform/Const"
	"TestPlatform/Struct"
	"TestPlatform/Util"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func ZdnsMakeLicense(Device Struct.MakeLicenseStruct) (string, error) {

	logrus := logrus.New()

	//校验机器是否可以链接
	client, conn_err := Device.SshStruct.Conn()
	if conn_err != nil {
		msg := "conn dev fail " + Device.SshStruct.Ipaddr + " " + conn_err.Error()
		return msg, errors.New(msg)
	}

	//校验机器是否安装了build
	exe_result, _ := Device.SshStruct.Execute(client, " lurker -h")
	if strings.Contains(exe_result, "No such file or directory") {
		msg := Device.SshStruct.Ipaddr + ":device not install lurker , please install ..."
		return msg, errors.New(msg)
	}

	//查看是否生成了machine.info
	exe_result, _ = Device.SshStruct.Execute(client, "ls /etc/machine.info")
	if strings.Contains(exe_result, "No such file or directory") {
		Device.SshStruct.Execute(client, "lurker -c fetch -o /etc/machine.info")
	}
	exe_result, exeErr := Device.SshStruct.Execute(client, "ls /etc/license.file")
	if exeErr == nil {
		logrus.Error(exeErr)
		msg := Device.SshStruct.Ipaddr + ":device exist /etc/license.file , exit ..."
		return msg, errors.New(msg)
	}

	programPath,_:= Util.GetPronPath()
	//制作license的本地文件夹
	tmpFilePath := programPath + Const.TempLicensePath + strings.Replace(Device.SshStruct.Ipaddr, ".", "_", -1) + "/"

	local_machine_info := tmpFilePath + "machine.info"
	local_pub_key := tmpFilePath + "pub.key"
	local_pri_key := tmpFilePath + "pri.key"
	local_licnese_file := tmpFilePath + "license.file"

	remote_machine_info := Const.RemoteLicnesePath + "machine.info"
	remote_pub_key := Const.RemoteLicnesePath + "pub.key"
	remote_pri_key := Const.RemoteLicnesePath + "pri.key"
	remote_licnese_file := Const.RemoteLicnesePath + "license.file"

	//测试代码
	Util.CreateDir(tmpFilePath)

	//Util.OsExecCmd("mkdir " + tmpFilePath)

	t := time.Now()
	time_now := fmt.Sprintf("%d-%d-%d %d:%d:%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	Device.SshStruct.Exec(client, "lurker -c fetch -o /etc/machine.info")
	Device.SshStruct.GetFile(local_machine_info, remote_machine_info, logrus)

	ddi_version, _ := Device.SshStruct.Execute(client, "rpm -qa | grep zddi")
	if Device.DdiVersion == "old" || strings.Contains(ddi_version, "3.10") || strings.Contains(ddi_version, "3.11") {
		cmds := fmt.Sprintf("license_old -p '%s' -q '%s' -i '%s' -l '%s' -e 'license_start_time:#:%s;license_valid_days:#:-1;license_remote_auth_flag:#:0;license_flag:#:1;root_change_password_flag:#:0;license_control_node_num:#:1;license_normal_flag:#:1;license_virtual_validate:#:1;user_name:#:test;license_type:#:0;device_type:#:0;elk:#:1;sgcloud_st_time:#:%s;sgcloud_vaild_days:#:-1;sgcloud_role:#:1' -r 'DNSv%s,DHCPv%s,ADDv%s,REGv0 0 0'", local_pub_key, local_pri_key, local_machine_info, local_licnese_file, time_now, time_now, Device.DnsVersion, Device.DhcpVersion, Device.AddVersion)
		Util.OsExecCmd(cmds)
	} else if Device.DdiVersion == "new" || strings.Contains(ddi_version, "3.13") || strings.Contains(ddi_version, "3.14") || strings.Contains(ddi_version, "3.15") || strings.Contains(ddi_version, "3.49") {
		cmds := fmt.Sprintf("license_new -p '%s' -q '%s' -i '%s' -l '%s' -e 'license_start_time:#:%s;license_valid_days:#:-1;license_remote_auth_flag:#:0;root_change_password_flag:#:0;device_type:#:0' -r 'DNSv%s,DHCPv%s,ADDv%s 0 0'", local_pub_key, local_pri_key, local_machine_info, local_licnese_file, time_now, Device.DnsVersion, Device.DhcpVersion, Device.AddVersion)
		Util.OsExecCmd(cmds)
	} else if len(Device.DdiVersion) != 0 && Device.DdiVersion != "new" && Device.DdiVersion != "old" {
		msg := Device.SshStruct.Ipaddr + ":device generate license fail , unknow version :" + Device.DdiVersion + " please input 'new'/'old', exit ..."
		return msg, errors.New(msg)
	} else {
		msg := Device.SshStruct.Ipaddr + ":device generate license fail , unknow version :" + ddi_version + " exit ..."
		return msg, errors.New(msg)
	}
	Device.SshStruct.PutFile(local_pub_key, remote_pub_key, logrus)
	Device.SshStruct.PutFile(local_pri_key, remote_pri_key, logrus)
	Device.SshStruct.PutFile(local_licnese_file, remote_licnese_file, logrus)
	Util.OsExecCmd("rm -rf " + tmpFilePath)
	if exe_result, _ = Device.SshStruct.Execute(client, "ls /etc/license.file"); strings.Contains(exe_result, "No such file or directory") {
		msg := Device.SshStruct.Ipaddr + ":device generate license fail , exit ..."
		return msg, errors.New(msg)
	} else {
		msg := Device.SshStruct.Ipaddr + ":device generate license succ , please use ..."
		return msg, errors.New(msg)
	}
}
