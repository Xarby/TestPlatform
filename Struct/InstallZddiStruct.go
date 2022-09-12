package Struct

import (
	"github.com/jlaffaye/ftp"
	"net"
	"strings"
	"time"
)

// 部署节点的结构

// FTP部署的结构
type FtpTask struct {
	Ftp         FtpStruct   `json:"ftp"`
	DnsVersion  int         `json:"dns_version"`
	AddVersion  int         `json:"add_version"`
	DhcpVersion int         `json:"dhcp_version"`
	ZddiPath    string      `json:"zddi_path"`
	BuildPath   string      `json:"build_path"`
	ZddiDevices []ScpStruct `json:"zddi_devices"`
}

func (ftp_task FtpTask) CheckFtpTask() (bool, string) {
	ip := net.ParseIP(ftp_task.Ftp.Ipaddr)
	if ip == nil {
		return false, ftp_task.Ftp.Ipaddr + " not lawful !"
	}
	if ftp_task.DnsVersion > 4 || ftp_task.DnsVersion < 0 {
		return false, "dns version input limit 0-3"
	}
	if ftp_task.AddVersion > 3 || ftp_task.AddVersion < 0 {
		return false, "dns version input limit 0-2"
	}
	if ftp_task.DhcpVersion > 2 || ftp_task.DhcpVersion < 0 {
		return false, "dns version input limit 0-1"
	}
	url := ftp_task.Ftp.Ipaddr + ":" + ftp_task.Ftp.Port
	ftp_client, conn_err := ftp.Dial(url, ftp.DialWithTimeout(5*time.Second))
	if conn_err != nil {
		return false, "conn ftp server" + ftp_task.Ftp.Ipaddr + " fail ..."
	} else {
		//假如连接成功
		login_err := ftp_client.Login(ftp_task.Ftp.User, ftp_task.Ftp.Password)

		if login_err != nil {
			return false, "login ftp server" + ftp_task.Ftp.Ipaddr + " fail ... " + ftp_task.Ftp.User + "/" + ftp_task.Ftp.Password
		} else {
			//假如登录成功
			_, read_build_err := ftp_client.FileSize(ftp_task.BuildPath)
			if read_build_err != nil {
				return false, "ftp server " + ftp_task.Ftp.Ipaddr + " not find file " + ftp_task.BuildPath
			}
			_, read_zddi_err := ftp_client.FileSize(ftp_task.ZddiPath)
			if read_zddi_err != nil {
				return false, "ftp server " + ftp_task.Ftp.Ipaddr + " not find file " + ftp_task.ZddiPath
			}
		}
	}
	for _, device := range ftp_task.ZddiDevices {
		_, conn_scp_err := device.Conn()
		if conn_scp_err != nil {
			return false, "conn fail zddi " + device.Ipaddr
		}
	}
	return true, "check succ start task !"
}

//  SCP部署的结构
type ScpTask struct {
	GetScpZddi  ScpStruct   `json:"get_scp_zddi"`
	GetScpBuild ScpStruct   `json:"get_scp_build"`
	DnsVersion  int         `json:"dns_version"`
	AddVersion  int         `json:"add_version"`
	DhcpVersion int         `json:"dhcp_version"`
	ZddiDevices []ScpStruct `json:"zddi_devices"`
}

func (scp_task ScpTask) CheckScpTask() (bool, string) {

	if scp_task.DnsVersion > 4 || scp_task.DnsVersion < 0 {
		return false, "dns version input limit 0-3"
	}
	if scp_task.AddVersion > 3 || scp_task.AddVersion < 0 {
		return false, "dns version input limit 0-2"
	}
	if scp_task.DhcpVersion > 2 || scp_task.DhcpVersion < 0 {
		return false, "dns version input limit 0-1"
	}

	//检查远端是否有build包

	if check_flag,msg :=check_scp_file_exist(scp_task.GetScpBuild);check_flag==false{
		return check_flag,msg
	}
	if check_flag,msg :=check_scp_file_exist(scp_task.GetScpZddi);check_flag==false{
		return check_flag,msg
	}

	for _, device := range scp_task.ZddiDevices {
		_, conn_scp_err := device.Conn()
		if conn_scp_err != nil {
			return false, "conn fail zddi " + device.Ipaddr
		}
	}
	return true, "check succ start task !"
}

func check_scp_file_exist(scp_dev ScpStruct) (bool, string) {
	scp_conn, conn_scp_err := scp_dev.Conn()
	defer scp_conn.Close()
	if conn_scp_err != nil {
		return false, "conn fail zddi " + scp_dev.Ipaddr
	} else {
		exe_result, err := scp_dev.Exec(scp_conn, "ls "+scp_dev.Path)
		if err != nil {
			return false, "exe get remote file cmd fail :" + scp_dev.Ipaddr + ":" + scp_dev.Path
		} else if strings.Contains(exe_result, "No such file or directory") {
			return false, "get remote build fail :" + scp_dev.Ipaddr + ":" + scp_dev.Path
		}
	}
	return true,"file"+scp_dev.Path+" exist"
}
