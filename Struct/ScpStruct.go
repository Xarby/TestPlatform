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
	sshclient *ssh.Client
}

func (scp_dev *ScpStruct) PutFile(local string, remote string) error {
	ssh_client, conn_err := scp_dev.Conn()
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
		logrus.Error("not exist file " + local)
	}
	logrus.Info("start push local file： " + local + " to " + scp_dev.Ipaddr + " " + remote + " ...")
	if push_err := scpclient.CopyFromFile(context.Background(), *f, remote, "0655"); push_err != nil {
		logrus.Error(push_err)
	}
	logrus.Info("push ok ...")
	return nil
}

func (scp_dev *SshStruct) GetFile(local string, remote string) (string, error) {

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

//zddi_path string, build_path string, dns_version string, add_version string, dhcp_version string, role string
func (dev *ScpStruct) InstallZddi(zddi_path string, build_path string, dns_version int, add_version int, dhcp_version int, role string) error {
	if strings.Contains(build_path,"3.15") || strings.Contains(build_path,"3.16"){
		dev.CheckBackup()
	}

	localZddiFileName := Const.ZddiFileMenuName + zddi_path
	localBuildFileName := Const.ZddiFileMenuName + build_path
	remoteZddiFileName := Const.ZddiRemoteFilePath + zddi_path
	remoteBuildFileName := Const.ZddiRemoteFilePath + build_path
	localPrivate, remotePrivate := Util.ChooseVersion(build_path)
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

	logrus.Info("-------------POST SCP FILE MD5 INFO-------------")
	logrus.Info("remote ipaddr: " + dev.Ipaddr)
	logrus.Info("local zddi md5: " + local_zddi_md5)
	logrus.Info("remote zddi md5: " + remote_zddi_md5)
	logrus.Info("local build md5: " + local_build_md5)
	logrus.Info("remote build md5: " + remote_build_md5)
	logrus.Info("-------------------------------------------")

	if get_zddi_md5_err == nil {
		if remote_zddi_md5 != local_zddi_md5 {
			logrus.Warning("remote dev :" + dev.Ipaddr + " file " + Const.ZddiRemoteFilePath + zddi_path + " Md5 diff local")
			logrus.Warning("local :" + local_zddi_md5)
			logrus.Warning("remote :" + remote_zddi_md5)
			dev.PutFile(localZddiFileName, remoteZddiFileName)
		} else {
			logrus.Info("remote dev :" + dev.Ipaddr + " file " + Const.ZddiRemoteFilePath + zddi_path + " Md5 same local ,skip push file")
		}

	} else {
		dev.PutFile(localZddiFileName, remoteZddiFileName)
	}

	if get_build_md5_err == nil {
		if remote_build_md5 != local_build_md5 {
			logrus.Warning("remote dev :" + dev.Ipaddr + " file " + Const.ZddiRemoteFilePath + build_path + " Md5 diff local")
			logrus.Warning("local :" + local_build_md5)
			logrus.Warning("remote:" + remote_build_md5)
			dev.PutFile(localBuildFileName, remoteBuildFileName)
		} else {
			logrus.Info("remote dev :" + dev.Ipaddr + " file " + Const.ZddiRemoteFilePath + build_path + " Md5 same local ,skip push file")
		}
	} else {
		dev.PutFile(localBuildFileName, remoteBuildFileName)
	}
	remote_private_md5, get_private_err := dev.Exec(ssh_client, "md5sum "+remotePrivate)
	local_private_md5, _ := Util.GetFileMd5(localPrivate)

	if get_private_err == nil {
		if strings.Contains(remote_private_md5, local_private_md5) == false {
			dev.PutFile(localPrivate, remotePrivate)
		} else {
			logrus.Error("remote exist file" + remotePrivate + " md5:" + remote_private_md5)
		}
	} else {
		dev.PutFile(localPrivate, remotePrivate)
	}

	logrus.Info("start install zddi " + dev.Ipaddr)
	dev.Exec(ssh_client, "rpm -ivh "+build_path+" --force")
	//dev.Exec(ssh_client,`echo  -n "" > /usr/local/appsys/normal/version/os_version`)
	//dev.Exec(ssh_client,`echo  '{"current_version":"bugfix-3.15.3.7","e_current_version":"bugfix-3.15.3.7"}' > /usr/local/appsys/normal/version/os_version`)
	if strings.Contains(zddi_path, "rpm") {
		dev.Exec(ssh_client, "rpm -ivh "+zddi_path+" --force")

	} else {
		dev.Exec(ssh_client, "tar -xzvf "+zddi_path)
	}
	dev.Exec(ssh_client, "tar -xzvf "+remotePrivate)
	dev.Exec(ssh_client, "cd /root/zddi-private/publisher/ && sh install")
	dev.Exec(ssh_client, "lurker -c fetch -o /etc/machine.info")

	if strings.Contains(build_path, "3.10") || strings.Contains(build_path, "3.11") {
		dev.Exec(ssh_client, fmt.Sprintf("publisher -c license -p /etc/pub.key -q /etc/pri.key -i /etc/machine.info -l /etc/license.file -N %d -H %d -A %d -R 1 -m 0", dns_version, dhcp_version, add_version))
		dev.Exec(ssh_client, fmt.Sprintf("sed -i 's#^lurker.*#lurker -c run -d /root/zddi -r %s -p /etc/pub.key -l /etc/license.file#g'  /etc/rc.d/rc.local", role))
		dev.Exec(ssh_client, "nohup /etc/rc.local &")
	} else if strings.Contains(build_path, "3.13") || strings.Contains(build_path, "3.14") {
		dev.Exec(ssh_client, fmt.Sprintf("publisher -c license -p /etc/pub.key -q /etc/pri.key -i /etc/machine.info -l /etc/license.file -N %d -H %d -A %d  -m 0", dns_version, dhcp_version, add_version))
		dev.Exec(ssh_client, fmt.Sprintf("sed -i 's#^lurker.*#lurker -c run -d /root/zddi -r %s -p /etc/pub.key -l /etc/license.file#g'  /etc/rc.d/rc.local", role))
		dev.Exec(ssh_client, "nohup /etc/rc.local &")
	} else if strings.Contains(build_path, "3.15") || strings.Contains(build_path, "3.16") || strings.Contains(build_path, "3.17") {
		dev.Exec(ssh_client, fmt.Sprintf("publisher -c license -p /etc/pub.key -q /etc/pri.key -i /etc/machine.info -l /etc/license.file -N %d -H %d -A %d -m 0", dns_version, dhcp_version, add_version))
		dev.Exec(ssh_client, "sed -i '/release/d' /etc/rc.d/rc.local")
		dev.Exec(ssh_client, fmt.Sprintf("sed -i '/CLISH_PATH/a /usr/local/appsys/normal/package/zdns_startmgr/zdns_startmgr_ctl release /root/zddi %s' /etc/rc.d/rc.local", role))
		dev.Exec(ssh_client, "nohup /etc/rc.local &")
		dev.Exec(ssh_client, "sed -i '/PermitRootLogin/d' /etc/ssh/sshd_config")
		dev.Exec(ssh_client, "sed -i '/StrictModes/a PermitRootLogin yes' /etc/ssh/sshd_config")
		dev.Exec(ssh_client, "service  sshd restart")
		//3.10||3.11||3.13||3.14版本处理方式
	}
	return nil
}
func (dev *SshStruct) CheckBackup() (string, error) {
	ssh_client, conn_err := dev.Conn()
	if conn_err != nil {
		return  "conn "+dev.Ipaddr+" fail", conn_err
	}
	defer ssh_client.Close()
	exec_zddi_result, _ := dev.Exec(ssh_client, "rpm -qa | grep zddi")
	//安装了rpm包
	if strings.Contains(exec_zddi_result, "zddi") {
		dev.Exec(ssh_client, "zdns-sys-recovery recovery")
		for i := 0; i < 30; i++ {
			port, _ := strconv.Atoi(dev.Port)
			conn_bool, conn_err := Util.TryConn(dev.Ipaddr, port)
			if conn_bool == true {
				logrus.Info(dev.Ipaddr + " recover succ！")
				return  dev.Ipaddr + " recover succ！", nil
			}
			logrus.Info(dev.Ipaddr+" recover later reboot fail！", conn_err)
			return  dev.Ipaddr + " recover later reboot fail！", conn_err
		}
		//未安装rpm包
	} else {
		scp := ScpStruct{
			SshStruct: *dev,
			Path:      "",
			Role:      "",
			sshclient: ssh_client,
		}
		scp.PutFile(Const.LocalRecoverToolPath, Const.RemoteRecoverToolPath)
		time.Sleep(time.Second * 5)
		dev.Exec(ssh_client, "tar -xzvf sys_recovery_test.tar.gz")
		dev.Exec(ssh_client, "sh "+Const.RecoverToolInstallPath)
		exe_result, _ := dev.Exec(ssh_client, "zdns-sys-recovery backup")
		//查看是否支持备份
		if strings.Contains(exe_result, "System backup is not supported") {
			logrus.Warning(dev.Ipaddr+" backup fail", errors.New(exe_result))
			return dev.Ipaddr + " backup fail", errors.New(exe_result)
		} else {
			logrus.Info(dev.Ipaddr+" backup succ", nil)
			return  dev.Ipaddr + " backup succ", nil
		}
		return  dev.Ipaddr + " backup succ", nil
	}
	return "", nil
}
