package Struct

import (
	"TestPlatform/Const"
	"TestPlatform/Util"
	"errors"
	"fmt"
	"github.com/bramvdbogaerde/go-scp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"
	"os"
	"strconv"
	"strings"
	"time"
)

type Devs struct {
	Devices []SshStruct `json:"devices"`
}

// SCP获取文件结构体
type ScpStruct struct {
	SshStruct
	Path      string `json:"path"`
	Role      string `json:"role"`
	RemoteIp  string `json:"remote_ip"`
	VIp       string `json:"vip"`
	sshclient *ssh.Client
}

func (ssh_dev *SshStruct) PutFile(local string, remote string, logrus *logrus.Logger) error {
	logrus.Debug("start push local file： " + local + " to " + ssh_dev.Ipaddr + " " + remote + " ...")
	ssh_client, conn_err := ssh_dev.Conn()
	if conn_err != nil {
		return conn_err
	}
	scpclient, new_client_err := scp.NewClientBySSH(ssh_client)
	defer scpclient.Close()
	if new_client_err != nil {
		logrus.Error("ssh change scp fail")
	}
	if _, file_err := os.Stat(local); file_err != nil {
		logrus.Error("local not exist file " + local)
	}
	f, file_err := os.Open(local)

	if file_err != nil {
		msg := "not exist file " + local
		logrus.Error(msg)
		return errors.New(msg)
	}
	logrus.Info("start push local file： " + local + " to " + ssh_dev.Ipaddr + " " + remote + " ...")
	if push_err := scpclient.CopyFromFile(context.Background(), *f, remote, "0655"); push_err != nil {
		logrus.Error(push_err)
	}
	logrus.Info("push ok ...")
	return nil
}

func (scp_dev SshStruct) GetFile(local string, remote string, logrus *logrus.Logger) (string, error) {

	//ssh连接
	ssh_client, conn_err := scp_dev.Conn()
	if conn_err != nil {
		return "conn " + scp_dev.Ipaddr + " error ", conn_err
	}
	//转scp
	scpclient, new_client_err := scp.NewClientBySSH(ssh_client)
	if new_client_err != nil {
		logrus.Error("ssh change scp fail")
	}
	defer scpclient.Close()
	//连接
	client, connErr := scp_dev.Conn()
	if connErr != nil {
		return "get", connErr
	}
	defer client.Close()

	//判断远程是否存在文件
	_, getErr := scp_dev.Exec(client, "ls "+remote)
	if getErr != nil {
		msg := scp_dev.Ipaddr + ":" + getErr.Error()
		return msg, errors.New(msg)
	}

	//下载文件
	srcFile, open_file_err := os.OpenFile(local, os.O_RDWR|os.O_CREATE, 0755)
	if open_file_err != nil {
		logrus.Error("open file " + local + " fail")
	}
	logrus.Info("start pull remote file： " + remote + " to " + local + " ...")
	if get_file_err := scpclient.CopyFromRemote(context.Background(), srcFile, remote); get_file_err != nil {
		logrus.Error("push err ...", get_file_err)
		return "get remote file :" + remote + "fail", get_file_err
	}

	return "get remote file" + remote + " succ...", nil
}

// zddi_path string, build_path string, dns_version string, add_version string, dhcp_version string, role string
func (dev *ScpStruct) InstallZddi(zddi_path string, build_path string, dns_version int, add_version int, dhcp_version int, role string, task_name string) error {

	logrus := Util.CreateLogger(Const.ZddiTaskLogPath+task_name, Const.ZddiTaskLogPath+task_name+"/"+dev.Ipaddr+".log")
	dev.InstallKey(logrus)
	//if _,check_err:=dev.CheckBackupTask();check_err!=nil{
	//	logrus.Warning(check_err.Error())
	//}else {
	//	dev.StartBackup()
	//}
	localZddiFileName := Const.ZddiFileMenuName + zddi_path
	localBuildFileName := Const.ZddiFileMenuName + build_path
	remoteZddiFileName := Const.ZddiRemoteFilePath + zddi_path
	remoteBuildFileName := Const.ZddiRemoteFilePath + build_path

	ssh_client, conn_err := dev.Conn()
	if conn_err != nil {
		return conn_err
	}
	defer ssh_client.Close()
	//传输文件

	remote_zddi_md5, get_zddi_md5_err := dev.Exec(ssh_client, "md5sum "+remoteZddiFileName+"| awk {'print $1'} ")
	local_zddi_md5, _ := Util.GetFileMd5(localZddiFileName)
	remote_build_md5, get_build_md5_err := dev.Exec(ssh_client, "md5sum "+remoteBuildFileName+"| awk {'print $1'} ")
	local_build_md5, _ := Util.GetFileMd5(localBuildFileName)

	logrus.Debug("remote ipaddr: " + dev.Ipaddr)
	logrus.Debug("local zddi md5: " + local_zddi_md5)
	logrus.Debug("remote zddi md5: " + remote_zddi_md5)
	logrus.Debug("local build md5: " + local_build_md5)
	logrus.Debug("remote build md5: " + remote_build_md5)

	if get_zddi_md5_err == nil {
		if remote_zddi_md5 != local_zddi_md5 {
			logrus.Warning("remote dev :" + dev.Ipaddr + " file " + Const.ZddiRemoteFilePath + zddi_path + " Md5 diff local")
			logrus.Warning("-------------------------------------------")
			logrus.Warning("local :" + local_zddi_md5)
			logrus.Warning("remote :" + remote_zddi_md5)
			logrus.Warning("-------------------------------------------")
			dev.PutFile(localZddiFileName, remoteZddiFileName, logrus)
		} else {
			logrus.Debug("remote dev :" + dev.Ipaddr + " file " + Const.ZddiRemoteFilePath + zddi_path + " Md5 same local ,skip push file")
		}

	} else {
		dev.PutFile(localZddiFileName, remoteZddiFileName, logrus)
	}

	if get_build_md5_err == nil {
		if remote_build_md5 != local_build_md5 {
			logrus.Warning("remote dev :" + dev.Ipaddr + " file " + Const.ZddiRemoteFilePath + build_path + " Md5 diff local")
			logrus.Warning("-------------------------------------------")
			logrus.Warning("local :" + local_build_md5)
			logrus.Warning("remote:" + remote_build_md5)
			logrus.Warning("-------------------------------------------")
			dev.PutFile(localBuildFileName, remoteBuildFileName, logrus)
		} else {
			logrus.Debug("remote dev :" + dev.Ipaddr + " file " + Const.ZddiRemoteFilePath + build_path + " Md5 same local ,skip push file")
		}
	} else {
		dev.PutFile(localBuildFileName, remoteBuildFileName, logrus)
	}

	logrus.Info("start install build pkg ...")
	dev.Exec(ssh_client, "rpm -ivh "+build_path)
	logrus.Info("end install build pkg ...")
	if strings.Contains(zddi_path, "rpm") {
		logrus.Debug("start install zddi pkg ...")
		dev.Exec(ssh_client, "rpm -ivh "+zddi_path)
		logrus.Info("end install zddi pkg ...")
	} else {
		logrus.Debug("start install zddi pkg ...")
		dev.Exec(ssh_client, "tar -xzvf "+zddi_path)
		logrus.Info("end install zddi pkg ...")
	}
	logrus.Info("start make license  ...")

	programPath, _ := Util.GetPronPath()
	//制作license的本地文件夹
	tmpFilePath := programPath + Const.TempLicensePath + strings.Replace(dev.Ipaddr, ".", "_", -1) + "/"
	//新建此文件
	Util.OsExecCmd("mkdir " + tmpFilePath)

	local_machine_info := tmpFilePath + "machine.info"
	local_pub_key := tmpFilePath + "pub.key"
	local_pri_key := tmpFilePath + "pri.key"
	local_licnese_file := tmpFilePath + "license.file"

	remote_machine_info := Const.RemoteLicnesePath + "machine.info"
	remote_pub_key := Const.RemoteLicnesePath + "pub.key"
	remote_pri_key := Const.RemoteLicnesePath + "pri.key"
	remote_licnese_file := Const.RemoteLicnesePath + "license.file"
	t := time.Now()
	time_now := fmt.Sprintf("%d-%d-%d %d:%d:%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	dev.Exec(ssh_client, "lurker -c fetch -o /etc/machine.info")
	if strings.Contains(build_path, "3.10") || strings.Contains(build_path, "3.11") {
		dev.GetFile(local_machine_info, remote_machine_info, logrus)
		cmds := fmt.Sprintf("license_old -p '%s' -q '%s' -i '%s' -l '%s' -e 'license_start_time:#:%s;license_valid_days:#:-1;license_remote_auth_flag:#:0;license_flag:#:1;root_change_password_flag:#:0;license_control_node_num:#:1;license_normal_flag:#:1;license_virtual_validate:#:1;user_name:#:test;license_type:#:0;device_type:#:0;elk:#:1;sgcloud_st_time:#:%s;sgcloud_vaild_days:#:-1;sgcloud_role:#:1' -r 'DNSv%s,DHCPv%s,ADDv%s,REGv0 0 0'", local_pub_key, local_pri_key, local_machine_info, local_licnese_file, time_now, time_now, strconv.Itoa(dns_version), strconv.Itoa(dhcp_version), strconv.Itoa(add_version))
		Util.OsExecCmd(cmds)
		dev.Exec(ssh_client, fmt.Sprintf("sed -i 's#^lurker.*#lurker -c run -d /root/zddi -r %s -p /etc/pub.key -l /etc/license.file#g'  /etc/rc.d/rc.local", Util.ChangeRole(role)))
	} else if strings.Contains(build_path, "3.13") || strings.Contains(build_path, "3.14") {
		dev.GetFile(local_machine_info, remote_machine_info, logrus)
		cmds := fmt.Sprintf("license_new -p '%s' -q '%s' -i '%s' -l '%s' -e 'license_start_time:#:%s;license_valid_days:#:-1;license_remote_auth_flag:#:0;root_change_password_flag:#:0;device_type:#:0' -r 'DNSv%s,DHCPv%s,ADDv%s 0 0'", local_pub_key, local_pri_key, local_machine_info, local_licnese_file, time_now, strconv.Itoa(dns_version), strconv.Itoa(dhcp_version), strconv.Itoa(add_version))
		Util.OsExecCmd(cmds)
		dev.Exec(ssh_client, fmt.Sprintf("sed -i 's#^lurker.*#lurker -c run -d /root/zddi -r %s -p /etc/pub.key -l /etc/license.file#g'  /etc/rc.d/rc.local", Util.ChangeRole(role)))
	} else if strings.Contains(build_path, "3.15") || strings.Contains(build_path, "3.16") || strings.Contains(build_path, "3.17") || strings.Contains(build_path, "3.49") {
		dev.Exec(ssh_client, "sed -i '/release/d' /etc/rc.d/rc.local")
		dev.Exec(ssh_client, fmt.Sprintf("sed -i '/CLISH_PATH/a /usr/local/appsys/normal/package/zdns_startmgr/zdns_startmgr_ctl release /root/zddi %s' /etc/rc.d/rc.local", Util.ChangeRole(role)))
		dev.Exec(ssh_client, "sed -i '/PermitRootLogin/d' /etc/ssh/sshd_config")
		dev.Exec(ssh_client, "sed -i '/StrictModes/a PermitRootLogin yes' /etc/ssh/sshd_config")
		dev.Exec(ssh_client, "service  sshd restart")
		dev.GetFile(local_machine_info, remote_machine_info, logrus)
		cmds := fmt.Sprintf("license_new -p '%s' -q '%s' -i '%s' -l '%s' -e 'license_start_time:#:%s;license_valid_days:#:-1;license_remote_auth_flag:#:0;root_change_password_flag:#:0;device_type:#:0' -r 'DNSv%s,DHCPv%s,ADDv%s 0 0'", local_pub_key, local_pri_key, local_machine_info, local_licnese_file, time_now, strconv.Itoa(dns_version), strconv.Itoa(dhcp_version), strconv.Itoa(add_version))
		Util.OsExecCmd(cmds)
	}
	//激活license
	dev.PutFile(local_pub_key, remote_pub_key, logrus)
	dev.PutFile(local_pri_key, remote_pri_key, logrus)
	dev.PutFile(local_licnese_file, remote_licnese_file, logrus)
	Util.OsExecCmd("rm -rf " + tmpFilePath)
	//启动服务
	dev.Exec(ssh_client, "nohup /etc/rc.local &")

	logrus.Info("start zddi server ...")
	logrus.Info(dev.Ipaddr + " install zddi succ , please wait server start .")
	return nil
}
func (dev SshStruct) CheckBackup() (string, error) {
	ssh_client, conn_err := dev.Conn()
	if conn_err != nil {
		return "conn " + dev.Ipaddr + " fail", conn_err
	}
	defer ssh_client.Close()
	return dev.StartBackup()
}

func (dev SshStruct) InstallKey(*logrus.Logger) (string, error) {
	ssh_client, conn_err := dev.Conn()
	if conn_err != nil {
		return "conn " + dev.Ipaddr + " fail", conn_err
	}
	defer ssh_client.Close()
	dev.Exec(ssh_client, "mkdir /root/.ssh")
	dev.Exec(ssh_client, "touch /root/.ssh/authorized_keys")
	dev.Exec(ssh_client, "echo '"+Const.SshKey+"' >> /root/.ssh/authorized_keys")
	dev.Exec(ssh_client, "chmod 600 /root/.ssh/authorized_keys")
	dev.Exec(ssh_client, "chmod 700 /root/.ssh")
	dev.Exec(ssh_client, "sed -i -e '/RSAAuthentication/d' /etc/ssh/sshd_config")
	dev.Exec(ssh_client, "sed -i -e '/PubkeyAuthentication/d' /etc/ssh/sshd_config")
	dev.Exec(ssh_client, "sed -i -e '/PermitRootLogin/d' /etc/ssh/sshd_config")
	dev.Exec(ssh_client, "echo 'RSAAuthentication yes' >> /etc/ssh/sshd_config")
	dev.Exec(ssh_client, "echo 'PubkeyAuthentication yes' >> /etc/ssh/sshd_config")
	dev.Exec(ssh_client, "echo 'PermitRootLogin yes' >> /etc/ssh/sshd_config")
	dev.Exec(ssh_client, "service sshd restart")
	logrus.Info("Install Key Succ")
	return "Install Key Succ", nil
}

func (device SshStruct) CheckBackupTask() (string, error) {
	logrus := Util.CreateLogger(Const.RecoveryLogPath, Const.RecoveryLogPath+"/recovery.log")

	//开始连接
	client, conn_err := device.Conn()
	if conn_err != nil {
		return "", errors.New("conn device " + device.Ipaddr + " fail")
	} else {
		logrus.Info("conn device " + device.Ipaddr + " succ !")
	}
	//校验磁盘空间
	logrus.Debug("start get disk info")
	use_disk, _ := device.Exec(client, "df -l / |tail -n 1  | awk '{print $3}'")
	available_disk, _ := device.Exec(client, "df -l / |tail -n 1  | awk '{print $4}'")
	use_disk_byte, _ := strconv.Atoi(strings.Replace(use_disk, "\n", "", -1))
	available_disk_byte, _ := strconv.Atoi(strings.Replace(available_disk, "\n", "", -1))
	use_disk_gib := use_disk_byte / 1048576
	available_gib := available_disk_byte / 1048576
	logrus.Debug("use:" + strconv.Itoa(use_disk_gib) + "gib")
	logrus.Debug("available:" + strconv.Itoa(available_gib) + "gib")
	if use_disk_gib > available_gib {
		err_msg := "check disk space fail , error:insufficient space: use space " + strconv.Itoa(use_disk_gib) + "GB" + "  available space " + strconv.Itoa(available_gib) + "GB"
		logrus.Error(err_msg)
		return "", errors.New(err_msg)
	} else {
		logrus.Info("check disk space succ , " + strconv.Itoa(use_disk_gib) + "GB" + "  available space " + strconv.Itoa(available_gib) + "GB")
	}
	device.InstallKey(logrus)

	file, process, rpm, _ := device.ShowStatu()
	if file == 0 && process == 0 && rpm == 0 {
		result_msg := device.Ipaddr + " is pure environment , start backup task , please see Status... "
		return result_msg, nil
	} else if file == 2 {
		result_msg := device.Ipaddr + " exist execute backup file , please use recovery , skip backup ..."
		logrus.Warning(result_msg)
		return result_msg, errors.New(result_msg)
	} else if process == 1 {
		result_msg := device.Ipaddr + " exist execute backup or recovery task , please wait ..."
		logrus.Warning(result_msg)
		return result_msg, errors.New(result_msg)
	} else if rpm == 0 {
		result_msg := device.Ipaddr + " environment exist zddi build pkg , please recovery environment ."
		logrus.Warning(result_msg)
		return result_msg, errors.New(result_msg)
	} else if rpm == 0 && file == 0 {
		result_msg := device.Ipaddr + " environment exist zddi build pkg and not find backup meunu, please formnt environment ."
		logrus.Warning(result_msg)
		return result_msg, errors.New(result_msg)
	} else {
		result_msg := device.Ipaddr + " environment exist zddi build pkg and not find backup meunu, please formnt environment ."
		logrus.Warning(result_msg)
		return result_msg, errors.New(result_msg)
	}
	return "", nil
}

func (device SshStruct) StartBackup() (string, error) {
	logrus := Util.CreateLogger(Const.RecoveryLogPath, Const.RecoveryLogPath+"/recovery.log")

	//开始连接
	client, conn_err := device.Conn()
	if conn_err != nil {
		return "", errors.New("conn device " + device.Ipaddr + " fail")
	} else {
		logrus.Info("conn device " + device.Ipaddr + " succ !")
	}
	//校验磁盘空间
	logrus.Debug("start get disk info")
	use_disk, _ := device.Exec(client, "df -l / |tail -n 1  | awk '{print $3}'")
	available_disk, _ := device.Exec(client, "df -l / |tail -n 1  | awk '{print $4}'")
	use_disk_byte, _ := strconv.Atoi(strings.Replace(use_disk, "\n", "", -1))
	available_disk_byte, _ := strconv.Atoi(strings.Replace(available_disk, "\n", "", -1))
	use_disk_gib := use_disk_byte / 1048576
	available_gib := available_disk_byte / 1048576
	logrus.Debug("use:" + strconv.Itoa(use_disk_gib) + "gib")
	logrus.Debug("available:" + strconv.Itoa(available_gib) + "gib")
	if use_disk_gib > available_gib {
		err_msg := "check disk space fail , error:insufficient space: use space " + strconv.Itoa(use_disk_gib) + "GB" + "  available space " + strconv.Itoa(available_gib) + "GB"
		logrus.Error(err_msg)
		return "", errors.New(err_msg)
	} else {
		logrus.Info("check disk space succ , " + strconv.Itoa(use_disk_gib) + "GB" + "  available space " + strconv.Itoa(available_gib) + "GB")
	}
	device.InstallKey(logrus)

	file, process, rpm, _ := device.ShowStatu()
	if file == 0 && process == 0 && rpm == 0 {
		result_msg := device.Ipaddr + " is pure environment , start backup task , please see Status... "

		if exe_result, _ := device.Exec(client, "rpm -qa | grep Rsync"); strings.Contains(exe_result, "Rsync") == false {
			warn_msg := "Environment not installed Rsync , please install Rsync ..."
			logrus.Warning(warn_msg)
			arch, _ := device.Execute(client, "arch")
			system_info, _ := device.Execute(client, "cat /etc/system-release")
			if strings.Contains(arch, "x86_64") && strings.Contains(system_info, "7.8") {
				device.PutFile(Const.RsyncFilePath+Const.RsyncFilePathCentos7X86, Const.ZddiRemoteFilePath+Const.RsyncFilePathCentos7X86, logrus)
				device.Execute(client, "rpm -ivh "+Const.RsyncFilePathCentos7X86)
			} else if strings.Contains(arch, "x86_64") && strings.Contains(system_info, "6.4") {
				device.PutFile(Const.RsyncFilePath+Const.RsyncFilePathCentos6X86, Const.ZddiRemoteFilePath+Const.RsyncFilePathCentos6X86, logrus)
				device.Execute(client, "rpm -ivh "+Const.RsyncFilePathCentos6X86)
			} else if strings.Contains(arch, "x86_64") && strings.Contains(system_info, "openEuler") {
				device.PutFile(Const.RsyncFilePath+Const.RsyncFilePathopenEulerX86, Const.ZddiRemoteFilePath+Const.RsyncFilePathopenEulerX86, logrus)
				device.Execute(client, "rpm -ivh "+Const.RsyncFilePathopenEulerX86)
			} else if strings.Contains(arch, "aarch64") && strings.Contains(system_info, "V10") {
				device.PutFile(Const.RsyncFilePath+Const.RsyncFilePathKylin10ARM, Const.ZddiRemoteFilePath+Const.RsyncFilePathKylin10ARM, logrus)
				device.Execute(client, "rpm -ivh "+Const.RsyncFilePathKylin10ARM)
			} else {
				err_msg := "Not super arch " + arch + " system info " + system_info
				logrus.Error(err_msg)
				return "", errors.New(err_msg)
			}

		} else {
			logrus.Info("device " + device.Ipaddr + " exist Rsync !")
		}
		device.PutFile(Const.LocalBackupShellFile, Const.RemoteBackupToolShellFile, logrus)
		device.Execute(client, "chmod 755 "+Const.RemoteBackupToolShellFile)
		device.Execute(client, "zdns-recovery-tool backup")
		logrus.Info(result_msg)
		return result_msg, nil
	} else if file == 2 {
		result_msg := device.Ipaddr + " exist execute backup file , please use recovery , skip backup ..."
		logrus.Warning(result_msg)
		return result_msg, nil
	} else if process == 1 {
		result_msg := device.Ipaddr + " exist execute backup or recovery task , please wait ..."
		logrus.Warning(result_msg)
		return result_msg, nil
	} else if rpm == 0 {
		result_msg := device.Ipaddr + " environment exist zddi build pkg , please recovery environment ."
		logrus.Warning(result_msg)
		return result_msg, nil
	} else if rpm == 0 && file == 0 {
		result_msg := device.Ipaddr + " environment exist zddi build pkg and not find backup meunu, please formnt environment ."
		logrus.Warning(result_msg)
		return result_msg, nil
	} else {
		result_msg := device.Ipaddr + " environment exist zddi build pkg and not find backup meunu, please formnt environment ."
		logrus.Warning(result_msg)
		return result_msg, nil
	}
	return "", nil
}

func (device SshStruct) ShowStatu() (int, int, int, error) {
	client, conn_err := device.Conn()
	logrus.Debug(device, "check backup status")
	if conn_err != nil {
		return 0, 0, 0, conn_err
	} else {
		backup_file := 0
		process_exist := 0
		zddi_rpm_exist := 0
		_, backup_file_exist_err := device.Execute(client, "ls "+Const.BackupWorkDir)
		process_exist_err, _ := device.Execute(client, "ps -ef | grep 'rsync'  | grep -v grep")
		zddi_rpm_exist_err, _ := device.Execute(client, "rpm -qa | grep zddi")
		if backup_file_exist_err == nil {
			if backup_final, _ := device.Execute(client, "cat /zdns_backup/run_log/sys_full_bak.log  | grep 'total size is'"); strings.Contains(backup_final, "total size is") {
				backup_file = 2
			} else {
				backup_file = 1
			}
		}
		if strings.Contains(process_exist_err, "zdns_backup") {
			process_exist = 1
		}
		if strings.Contains(zddi_rpm_exist_err, "zddi") {
			fmt.Println(zddi_rpm_exist_err)
			zddi_rpm_exist = 1
		}
		return backup_file, process_exist, zddi_rpm_exist, nil
	}
}

func (device SshStruct) ExistFileOrPath(client *ssh.Client, path string) bool {
	if exe_result, exe_error := device.Execute(client, "ls "+path); exe_error != nil {
		return false
	} else if strings.Contains(exe_result, "No such file or directory") {
		return false
	} else {
		return true
	}
}
