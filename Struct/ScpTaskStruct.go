package Struct

import (
	"TestPlatform/Util"
	"errors"
	"fmt"
	"log"
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
		return "dns version input limit 0-3", errors.New("dns version input limit 0-3")
	}
	if scp_task.AddVersion >= 3 || scp_task.AddVersion < 0 {
		return "add version input limit 0-2", errors.New("add version input limit 0-2")
	}
	if scp_task.DhcpVersion >= 2 || scp_task.DhcpVersion < 0 {
		return "dhcp version input limit 0-1", errors.New("dhcp version input limit 0-1")
	}
	//校验juese
	for _, scp_struct := range scp_task.ZddiDevices {
		if check_inio, check_err := check_task_role(scp_struct); check_err != nil {
			return check_inio, check_err
		}
	}
	if scp_task.Colony == true {
		var count_master int
		for _, dev := range scp_task.ZddiDevices {
			log.Println(dev.Role)
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
		return msg, check_err
	}
	if msg, check_err := check_scp_file_exist(scp_task.GetScpZddi); check_err != nil {
		return msg, check_err
	}
	for _, device := range scp_task.ZddiDevices {
		_, conn_scp_err := device.Conn()
		if conn_scp_err != nil {
			return "conn fail zddi " + device.Ipaddr, conn_scp_err
		}
	}
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

	var master_ip string
	var count = 30
	//编辑所有节点
	for _, scp_dev := range scp_task.ZddiDevices {
		fmt.Println(scp_dev.Ipaddr)
		var port int
		//得到master
		if scp_dev.Role == "master" {
			master_ip = scp_dev.Ipaddr
			port = 443
		} else {
			if strings.Contains(scp_task.GetScpBuild.Path,"3.15")||strings.Contains(scp_task.GetScpBuild.Path,"3.16")||strings.Contains(scp_task.GetScpBuild.Path,"3.17"){
				port = 4583
			}else {
				port = 20123
			}
		}
		//判断节点是否开始服务
		for {
			if conn_flag, conn_err := Util.TryConn(scp_dev.Ipaddr, port); conn_flag == true {
				break
			} else {
				count--
				if count <0 {
					log.Println(master_ip,errors.New("count exceed 30!"))
					return master_ip,errors.New("count exceed 30!")
				}
				log.Println(conn_err)
				log.Println("sleep 1s")
				time.Sleep(time.Second)
			}
		}
	}

	//获取master的IP  并循环等待master的443端口开放 slave的4583端口开放
	fmt.Println("all open port")

	fmt.Println(master_ip)
	//编辑master的本机IP
	time.Sleep(time.Second*5)
	if requestsErr := Util.PostRequests("PUT", fmt.Sprintf("https://%s:20120/groups/local/members/master", master_ip), []byte(fmt.Sprintf(`{
											"group": "local",
											"name": "master",
											"id": "master",
											"ip": "%s",
											"positionX": "0.423",
											"positionY": "0.406"
										}`, master_ip))); requestsErr != nil {
		log.Println("put master ip faile " + master_ip+requestsErr.Error())
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
