package Struct

import (
	"TestPlatform/Util"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

//  SCP部署的结构
type ScpTask struct {
	GetScpZddi  ScpStruct   `json:"get_scp_zddi"`
	GetScpBuild ScpStruct   `json:"get_scp_build"`
	DnsVersion  int         `json:"dns_version"`
	AddVersion  int         `json:"add_version"`
	DhcpVersion int         `json:"dhcp_version"`
	ZddiDevices []ScpStruct `json:"zddi_devices"`
	Colony      bool        `json:"colony"`
}

func (scp_task ScpTask) CheckScpTask() (string, error) {
	//校验版本
	if scp_task.DnsVersion >= 4 || scp_task.DnsVersion < 0 {
		logrus.Error("dns version input limit 0-3")
		return "dns version input limit 0-3", errors.New("dns version input limit 0-3")
	}
	if scp_task.AddVersion >= 3 || scp_task.AddVersion < 0 {
		logrus.Error("add version input limit 0-2")
		return "add version input limit 0-2", errors.New("add version input limit 0-2")
	}
	if scp_task.DhcpVersion >= 2 || scp_task.DhcpVersion < 0 {
		logrus.Error("dhcp version input limit 0-1")
		return "dhcp version input limit 0-1", errors.New("dhcp version input limit 0-1")
	}
	//校验角色
	for _, scp_struct := range scp_task.ZddiDevices {
		if check_inio, check_err := check_task_role(scp_struct); check_err != nil {
			return check_inio, check_err
		}
	}
	if scp_task.Colony == true {
		var count_master int
		for _, dev := range scp_task.ZddiDevices {
			if dev.Role == "master" {
				count_master++
			}
		}
		if count_master >= 2 {
			return "colony model only one master role!", errors.New("colony model only master role! :")
		}
	}
	//检查远端是否有build包
	if msg, check_err := check_scp_file_exist(scp_task.GetScpBuild); check_err != nil {
		logrus.Error(msg, check_err)
		return msg, check_err
	}
	if msg, check_err := check_scp_file_exist(scp_task.GetScpZddi); check_err != nil {
		logrus.Error(msg, check_err)
		return msg, check_err
	}
	for _, device := range scp_task.ZddiDevices {
		client, conn_scp_err := device.Conn()
		if conn_scp_err != nil {
			return "conn fail zddi " + device.Ipaddr, conn_scp_err
		} else if strings.Contains(scp_task.GetScpBuild.Path, "3.15") || strings.Contains(scp_task.GetScpBuild.Path, "3.16") {

			exe_rpm_result, _ := device.Exec(client, "rpm -qa | grep zddi")
			exe_recovery_result, _ := device.Exec(client, "zdns-sys-recovery")
			//假如安装了rpm包
			if strings.Contains(exe_rpm_result, "zddi") {
				//安装了工具未备份 则为污染
				if strings.Contains(exe_recovery_result, "command not found") {
					return "environmental pollution", errors.New("environmental pollution")
				//未进行备份 则为污染
				} else {
					exe_check_result, _ := device.Exec(client, "zdns-sys-recovery check")
					if strings.Contains(exe_check_result, "successfully") == false {
						logrus.Error(device.Ipaddr + " environmental pollution")
						return "environmental pollution", errors.New("environmental pollution")
					}
				}
			}
		}
	}
	logrus.Info("check scp succ start task !")
	return "check succ start task !", nil
}

//校验scp获取是否存在文件函数
func check_scp_file_exist(scp_dev ScpStruct) (string, error) {
	scp_conn, conn_scp_err := scp_dev.Conn()
	defer func() {
		if scp_conn != nil {
			scp_conn.Close()
		}
	}()
	if conn_scp_err != nil {
		return "conn fail zddi " + scp_dev.Ipaddr, conn_scp_err
	} else {
		exe_result, err := scp_dev.Exec(scp_conn, "ls "+scp_dev.Path)
		if err != nil {
			return "exe get remote file cmd fail :" + scp_dev.Ipaddr + ":" + scp_dev.Path, err
		} else if strings.Contains(exe_result, "No such file or directory") {
			return "get remote build fail :" + scp_dev.Ipaddr + ":" + scp_dev.Path, err
		}
	}
	return "file" + scp_dev.Path + " exist", nil
}

//校验角色函数
func check_task_role(scp_struct ScpStruct) (string, error) {
	if scp_struct.Role != "master" && scp_struct.Role != "slave" && scp_struct.Role != "backup" {
		return scp_struct.Ipaddr + " role: '" + scp_struct.Role + "' not in master/slave/backup", errors.New("not in master/slave/backup")
	} else {
		return "check role succ:" + scp_struct.Ipaddr + " role: '" + scp_struct.Role, nil
	}
}
func (scp_task ScpTask) ScpStartCreateColony() (string, error) {
	logrus.Info("start Create Colony!")
	var master_ip string
	//获取master的IP  并循环等待master的443端口开放 slave的4583端口开放
	for _, scp_dev := range scp_task.ZddiDevices {
		var port int
		if scp_dev.Role == "master" {
			master_ip = scp_dev.Ipaddr
			port = 443
		} else {
			if strings.Contains(scp_task.GetScpBuild.Path, "3.15") || strings.Contains(scp_task.GetScpBuild.Path, "3.16") || strings.Contains(scp_task.GetScpBuild.Path, "3.17") {
				port = 4583
			} else {
				port = 20123
			}
		}
		var count = 0
		for {
			if conn_flag, conn_err := Util.TryConn(scp_dev.Ipaddr, port); conn_flag == true {
				logrus.Info("conn ip:" + scp_dev.Ipaddr + " port" + strconv.Itoa(port) + "succ")
				break
			} else {
				//时间超过 60则放弃
				if count > 60 {
					return "", errors.New("start server time exceed 60s")
				} else {
					//时间+3秒并打印
					count = count + 3
					logrus.Warning(scp_dev.Ipaddr+" sleep 3s ", conn_err)
					time.Sleep(time.Second * 3)
				}
			}
		}
	}
	logrus.Info("check all dev port open!")
	logrus.Info("master_ip:" + master_ip)
	//编辑master的本机IP
	logrus.Info("start put master cloud center ip!")
	time.Sleep(time.Second * 5)
	if requestsErr := Util.PostRequests("PUT", fmt.Sprintf("https://%s:20120/groups/local/members/master", master_ip), []byte(fmt.Sprintf(`{
											"group": "local",
											"name": "master",
											"id": "master",
											"ip": "%s",
											"positionX": "0.423",
											"positionY": "0.406"
										}`, master_ip))); requestsErr != nil {
		logrus.Error("put master ip faile " + master_ip + requestsErr.Error())
		return "put master ip faile " + master_ip, requestsErr
	}
	//遍历非master添加节点
	for _, scp_dev := range scp_task.ZddiDevices {
		if scp_dev.Ipaddr == master_ip {
			continue
		} else {
			//判断版本是否为3.15之后
			name := scp_dev.Role + scp_dev.Ipaddr[strings.LastIndex(scp_dev.Ipaddr, ".")+1:]
			if strings.Contains(scp_task.GetScpBuild.Path, "3.15") || strings.Contains(scp_task.GetScpBuild.Path, "3.16") || strings.Contains(scp_task.GetScpBuild.Path, "3.17") {
				body := []byte(fmt.Sprintf(`{
					"name": "%s",
					"ip": "%s",
					"username": "admin",
					"password": "admincns",
					"role": "%s",
					"group": "local",
					"is_extend":"no"}`, name, scp_dev.Ipaddr, scp_dev.Role))
				if requestsErr := Util.PostRequests("POST", fmt.Sprintf("https://%s:20120/groups/local/members", master_ip), body); requestsErr != nil {
					logrus.Error("add node"+scp_dev.Ipaddr+"role:"+scp_dev.Role+"fail", requestsErr)
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
					logrus.Error("add node"+scp_dev.Ipaddr+"role:"+scp_dev.Role+"fail", requestsErr)
					return "add node" + scp_dev.Ipaddr + "role:" + scp_dev.Role + "fail", requestsErr
				}
			}
		}
	}
	logrus.Info("add all node succ task shutdown !!!")
	return "add all node succ", nil
}
