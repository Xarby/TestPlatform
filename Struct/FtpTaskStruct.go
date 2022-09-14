package Struct

import (
	"TestPlatform/Util"
	"errors"
	"fmt"
	"github.com/jlaffaye/ftp"
	"log"
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
	Colony      bool        `json:"colony"`
}

func (ftp_task FtpTask) CheckFtpTask() (string, error) {
	ip := net.ParseIP(ftp_task.Ftp.Ipaddr)
	if ip == nil {
		return ftp_task.Ftp.Ipaddr + " illegal!", errors.New(ftp_task.Ftp.Ipaddr + " illegal!")
	}
	if ftp_task.DnsVersion >= 4 || ftp_task.DnsVersion < 0 {
		return "dns version input limit 0-3", errors.New("dns version input limit 0-3")
	}
	if ftp_task.AddVersion >= 3 || ftp_task.AddVersion < 0 {
		return "add version input limit 0-2", errors.New("add version input limit 0-2")
	}
	if ftp_task.DhcpVersion >= 2 || ftp_task.DhcpVersion < 0 {
		return "dhcp version input limit 0-1", errors.New("dhcp version input limit 0-1")
	}

	//校验角色信息是否出现异常
	for _, scp_struct := range ftp_task.ZddiDevices {
		if check_inio, check_err := check_task_role(scp_struct); check_err != nil {
			return check_inio, check_err
		}
	}
	//校验集群模式下是否存在两个master
	if ftp_task.Colony == true {
		var count_master int
		for _, dev := range ftp_task.ZddiDevices {
			log.Println(dev.Role)
			if dev.Role == "master" {
				count_master++
			}
		}
		if count_master >= 2 {
			return "colony model only one master role!", errors.New("colony model only master role! :")
		}
	}
	url := ftp_task.Ftp.Ipaddr + ":" + ftp_task.Ftp.Port
	ftp_client, conn_err := ftp.Dial(url, ftp.DialWithTimeout(5*time.Second))
	if conn_err != nil {
		return "conn ftp server" + ftp_task.Ftp.Ipaddr + " fail ...", conn_err
	} else {
		//假如连接成功
		login_err := ftp_client.Login(ftp_task.Ftp.User, ftp_task.Ftp.Password)
		if login_err != nil {
			return "login ftp server" + ftp_task.Ftp.Ipaddr + " fail ... " + ftp_task.Ftp.User + "/" + ftp_task.Ftp.Password, login_err
		} else {
			//假如登录成功
			_, read_build_err := ftp_client.FileSize(ftp_task.BuildPath)
			if read_build_err != nil {
				return "ftp server " + ftp_task.Ftp.Ipaddr + " not find file " + ftp_task.BuildPath, read_build_err
			}
			_, read_zddi_err := ftp_client.FileSize(ftp_task.ZddiPath)
			if read_zddi_err != nil {
				return "ftp server " + ftp_task.Ftp.Ipaddr + " not find file " + ftp_task.ZddiPath, read_zddi_err
			}
		}
	}
	for _, device := range ftp_task.ZddiDevices {
		_, conn_scp_err := device.Conn()
		if conn_scp_err != nil {
			return "conn fail zddi " + device.Ipaddr, conn_scp_err
		}
	}
	return "check succ start task !", nil
}

func (ftp_task FtpTask) FtpStartCreateColony() (string, error) {
	var master_ip string

	//获取master的IP  并循环等待master的443端口开放 slave的4583端口开放
	for _, scp_dev := range ftp_task.ZddiDevices {
		fmt.Println(scp_dev.Ipaddr)
		var port int
		if scp_dev.Role == "master" {
			master_ip = scp_dev.Ipaddr
			port = 443
		} else {
			if strings.Contains(ftp_task.BuildPath, "3.15") || strings.Contains(ftp_task.BuildPath, "3.16") || strings.Contains(ftp_task.BuildPath, "3.17") {
				port = 4583
			} else {
				port = 20123
			}
		}
		for {
			if conn_flag, conn_err := Util.TryConn(scp_dev.Ipaddr, port); conn_flag == true {
				break
			} else {
				log.Println(conn_err)
				log.Println("sleep 1s")
				time.Sleep(time.Second)
			}
		}
	}
	fmt.Println("all open port")

	fmt.Println(master_ip)
	//编辑master的本机IP
	time.Sleep(time.Second * 5)
	if requestsErr := Util.PostRequests("PUT", fmt.Sprintf("https://%s:20120/groups/local/members/master", master_ip), []byte(fmt.Sprintf(`{
											"group": "local",
											"name": "master",
											"id": "master",
											"ip": "%s",
											"positionX": "0.423",
											"positionY": "0.406"
										}`, master_ip))); requestsErr != nil {
		log.Println("put master ip faile " + master_ip + requestsErr.Error())
		return "put master ip faile " + master_ip, requestsErr
	}
	//遍历非master添加节点
	for _, scp_dev := range ftp_task.ZddiDevices {
		if scp_dev.Ipaddr == master_ip {
			continue
		} else {
			//判断版本是否为3.15之后
			name := scp_dev.Role + scp_dev.Ipaddr[strings.LastIndex(scp_dev.Ipaddr, ".")+1:]
			if strings.Contains(ftp_task.BuildPath, "3.15") || strings.Contains(ftp_task.BuildPath, "3.16") || strings.Contains(ftp_task.BuildPath, "3.17") {
				body := []byte(fmt.Sprintf(`{
					"name": "%s",
					"ip": "%s",
					"username": "admin",
					"password": "admincns",
					"role": "%s",
					"group": "local",
					"is_extend":"no"}`, name, scp_dev.Ipaddr, scp_dev.Role))
				if requestsErr := Util.PostRequests("POST", fmt.Sprintf("https://%s:20120/groups/local/members", master_ip), body); requestsErr != nil {
					return "add node" + scp_dev.Ipaddr + "role:" + scp_dev.Role + "fail", requestsErr
				}
			} else {
				body := []byte(fmt.Sprintf(`{
					"name": "%s",
					"ip": "%s",
					"username": "admin",
					"password": "admincns",
					"role": "%s",
					"group": "local"}`, name, scp_dev.Ipaddr, scp_dev.Role))
				if requestsErr := Util.PostRequests("POST", fmt.Sprintf("https://%s:20120/groups/local/members", master_ip), body); requestsErr != nil {
					return "add node" + scp_dev.Ipaddr + "role:" + scp_dev.Role + "fail", requestsErr
				}
			}
		}
	}
	return "add all node succ", nil
}
